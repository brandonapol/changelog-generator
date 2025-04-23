package main

import (
	"fmt"
	"os"

	"changelog-generator/internal/changelog"
)

func main() {
	// Prompt user for the list of repositories
	repositories, err := changelog.PromptForRepositories()
	if err != nil {
		fmt.Printf("Error prompting for repositories: %v\n", err)
		os.Exit(1)
	}

	for _, repo := range repositories {
		fmt.Printf("Generating changelog for repository: %s\n", repo)

		// Fetch all tags in the repository
		tags, err := changelog.FetchTags(repo)
		if err != nil {
			fmt.Printf("Error fetching tags for repository %s: %v\n", repo, err)
			continue
		}

		// Prompt user if they want to create a new tag
		createNewTag, newTag, err := changelog.PromptForNewTag(tags[len(tags)-1])
		if err != nil {
			fmt.Printf("Error prompting for new tag for repository %s: %v\n", repo, err)
			continue
		}
		if createNewTag {
			if err := changelog.CreateTag(repo, newTag); err != nil {
				fmt.Printf("Error creating new tag for repository %s: %v\n", repo, err)
				continue
			}
			tags = append(tags, newTag)
		}

		// Prompt user to select from and to tags
		fromTag, toTag, err := changelog.PromptForTags(repo, tags)
		if err != nil {
			fmt.Printf("Error selecting tags for repository %s: %v\n", repo, err)
			continue
		}

		// Fetch commits between the selected tags
		commits, err := changelog.FetchCommits(repo, fromTag, toTag)
		if err != nil {
			fmt.Printf("Error fetching commits for repository %s: %v\n", repo, err)
			continue
		}

		// Prompt user to select commits to include
		selectedCommits, err := changelog.PromptForCommits(commits)
		if err != nil {
			fmt.Printf("Error selecting commits for repository %s: %v\n", repo, err)
			continue
		}

		// Generate the changelog content
		changelogContent := changelog.GenerateChangelog(selectedCommits, fromTag, toTag)

		// Render the changelog in Markdown format
		markdownChangelog, err := changelog.RenderMarkdown(changelogContent)
		if err != nil {
			fmt.Printf("Error rendering Markdown for repository %s: %v\n", repo, err)
			continue
		}

		// Render the changelog in HTML format
		htmlChangelog, err := changelog.RenderHTML(changelogContent)
		if err != nil {
			fmt.Printf("Error rendering HTML for repository %s: %v\n", repo, err)
			continue
		}

		// Write the Markdown changelog to internal/changelog/output/CHANGELOG.md
		if err := os.WriteFile("internal/changelog/output/CHANGELOG.md", []byte(markdownChangelog), 0644); err != nil {
			fmt.Printf("Error writing CHANGELOG.md for repository %s: %v\n", repo, err)
			continue
		}

		// Write the HTML changelog to internal/changelog/output/release-notes.html
		if err := os.WriteFile("internal/changelog/output/release-notes.html", []byte(htmlChangelog), 0644); err != nil {
			fmt.Printf("Error writing release-notes.html for repository %s: %v\n", repo, err)
			continue
		}

		fmt.Printf("Changelog generated successfully for repository: %s\n", repo)
	}
}
