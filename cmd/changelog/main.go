package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brandonapol/changelog-generator/internal/changelog"
)

// processRepository handles all steps for a single repository
func processRepository(repo string) error {
	fmt.Printf("Generating changelog for repository: %s\n", repo)

	// Fetch all tags in the repository
	tags, err := changelog.FetchTags(repo)
	if err != nil {
		return fmt.Errorf("error fetching tags: %w", err)
	}

	// Handle tag operations
	tags, err = handleTagOperations(repo, tags)
	if err != nil {
		return err
	}

	// Get selected tags
	fromTag, toTag, err := changelog.PromptForTags(repo, tags)
	if err != nil {
		return fmt.Errorf("error selecting tags: %w", err)
	}

	// Get commits and generate changelog
	changelogContent, err := generateChangelogContent(repo, fromTag, toTag)
	if err != nil {
		return err
	}

	// Write output files
	if err := writeOutputFiles(changelogContent); err != nil {
		return err
	}

	fmt.Printf("Changelog generated successfully for repository: %s\n", repo)
	return nil
}

// handleTagOperations handles operations related to tag creation
func handleTagOperations(repo string, tags []string) ([]string, error) {
	createNewTag, newTag, err := changelog.PromptForNewTag(tags[len(tags)-1])
	if err != nil {
		return nil, fmt.Errorf("error prompting for new tag: %w", err)
	}

	if createNewTag {
		if err := changelog.CreateTag(repo, newTag); err != nil {
			return nil, fmt.Errorf("error creating new tag: %w", err)
		}
		return append(tags, newTag), nil
	}

	return tags, nil
}

// generateChangelogContent fetches commits and generates changelog content
func generateChangelogContent(repo, fromTag, toTag string) (string, error) {
	// Fetch commits between the selected tags
	commits, err := changelog.FetchCommits(repo, fromTag, toTag)
	if err != nil {
		return "", fmt.Errorf("error fetching commits: %w", err)
	}

	// Prompt user to select commits to include
	selectedCommits, err := changelog.PromptForCommits(commits)
	if err != nil {
		return "", fmt.Errorf("error selecting commits: %w", err)
	}

	// Create a map of repository paths to selected commits
	repoCommits := map[string][]string{
		repo: selectedCommits,
	}

	// Generate the changelog content
	return changelog.GenerateChangelog(repoCommits), nil
}

// writeOutputFiles writes the generated changelog to markdown and HTML files
func writeOutputFiles(changelogContent string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll("internal/changelog/output", 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Render and write Markdown
	if err := changelog.RenderMarkdown(changelogContent); err != nil {
		return fmt.Errorf("error rendering markdown: %w", err)
	}

	// Check if release-notes.html exists
	if _, err := os.Stat("release-notes.html"); err == nil {
		// File exists, render and write HTML
		if err := changelog.RenderHTML(changelogContent); err != nil {
			return fmt.Errorf("error rendering HTML: %w", err)
		}
	} else if os.IsNotExist(err) {
		// File doesn't exist, skip HTML rendering
		fmt.Println("release-notes.html not found, skipping HTML generation")
	} else {
		// Some other error occurred
		return fmt.Errorf("error checking release-notes.html: %w", err)
	}

	return nil
}

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Prompt user for the list of repositories
	repositories, err := changelog.PromptForRepositories()
	if err != nil {
		fmt.Printf("Error prompting for repositories: %v\n", err)
		os.Exit(1)
	}

	// Process each repository
	for _, repo := range repositories {
		// Change into the current working directory
		if err := os.Chdir(cwd); err != nil {
			fmt.Printf("Error changing directory: %v\n", err)
			continue
		}

		if err := processRepository(repo); err != nil {
			fmt.Printf("Error processing repository %s: %v\n", repo, err)
		}
	}

	// Change back to the original directory
	if err := os.Chdir(filepath.Dir(os.Args[0])); err != nil {
		fmt.Printf("Error changing back to original directory: %v\n", err)
	}
}
