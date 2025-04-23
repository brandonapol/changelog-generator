package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Extract app version from the 'to' tag
	appVersion := toTag

	// Generate release date
	releaseDate := time.Now().Format("January 2, 2006")

	// Get commits and generate changelog
	changelogContent, features, bugfixes, others, err := generateChangelogContent(repo, fromTag, toTag)
	if err != nil {
		return err
	}

	// Search for CHANGELOG.md and release-notes.html files
	var changelogFile, releaseNotesFile string
	err = filepath.Walk(repo, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			if info.Name() == "CHANGELOG.md" {
				changelogFile = path
			} else if info.Name() == "release-notes.html" {
				releaseNotesFile = path
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error searching for files: %w", err)
	}

	// Write output files
	if err := writeOutputFiles(changelogContent, appVersion, releaseDate, features, bugfixes, others, changelogFile, releaseNotesFile); err != nil {
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

	// Return original tags if user opted not to create a new one
	return tags, nil
}

// generateChangelogContent fetches commits and generates changelog content
func generateChangelogContent(repo, fromTag, toTag string) (string, []string, []string, []string, error) {
	// Fetch commits between the selected tags
	commits, err := changelog.FetchCommits(repo, fromTag, toTag)
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error fetching commits: %w", err)
	}

	// Prompt user to select commits to include
	selectedCommits, err := changelog.PromptForCommits(commits)
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error selecting commits: %w", err)
	}

	// Create a map of repository paths to selected commits
	repoCommits := map[string][]string{
		repo: selectedCommits,
	}

	// Generate the changelog content
	changelogContent := changelog.GenerateChangelog(repoCommits)

	// Categorize the selected commits
	var features, bugfixes, others []string
	for _, commit := range selectedCommits {
		if strings.HasPrefix(commit, "feat:") {
			features = append(features, strings.TrimPrefix(commit, "feat: "))
		} else if strings.HasPrefix(commit, "fix:") {
			bugfixes = append(bugfixes, strings.TrimPrefix(commit, "fix: "))
		} else {
			others = append(others, commit)
		}
	}

	return changelogContent, features, bugfixes, others, nil
}

// writeOutputFiles writes the generated changelog to markdown and HTML files
func writeOutputFiles(changelogContent string, appVersion, releaseDate string, features, bugfixes, others []string, changelogFile, releaseNotesFile string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll("internal/changelog/output", 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Render and write Markdown
	fmt.Printf("Rendering Markdown changelog to %s...\n", changelogFile)
	if err := changelog.RenderMarkdown(changelogContent, appVersion, releaseDate, changelogFile); err != nil {
		return fmt.Errorf("error rendering markdown: %w", err)
	}
	fmt.Printf("Successfully rendered Markdown changelog to %s\n", changelogFile)

	// Render and write HTML
	fmt.Printf("Rendering HTML changelog to %s...\n", releaseNotesFile)
	if err := changelog.RenderHTML(changelogContent, appVersion, releaseDate, features, bugfixes, others, releaseNotesFile); err != nil {
		return fmt.Errorf("error rendering HTML: %w", err)
	}
	fmt.Printf("Successfully rendered HTML changelog to %s\n", releaseNotesFile)

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
