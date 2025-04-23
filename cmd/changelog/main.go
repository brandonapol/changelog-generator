package main

import (
	"flag"
	"fmt"
	"os"

	"changelog-generator/internal/changelog"
)

var (
	repoFlag = flag.String("repo", "", "Path to the Git repository")
)

func main() {
	flag.Parse()

	if *repoFlag == "" {
		fmt.Println("Please provide the path to the Git repository using the -repo flag")
		os.Exit(1)
	}

	// Fetch all tags in the repository
	tags, err := changelog.FetchTags(*repoFlag)
	if err != nil {
		fmt.Printf("Error fetching tags: %v\n", err)
		os.Exit(1)
	}

	// Prompt user if they want to create a new tag
	createNewTag, newTag, err := changelog.PromptForNewTag(tags[len(tags)-1])
	if err != nil {
		fmt.Printf("Error prompting for new tag: %v\n", err)
		os.Exit(1)
	}
	if createNewTag {
		if err := changelog.CreateTag(*repoFlag, newTag); err != nil {
			fmt.Printf("Error creating new tag: %v\n", err)
			os.Exit(1)
		}
		tags = append(tags, newTag)
	}

	// Prompt user to select from and to tags
	fromTag, toTag, err := changelog.PromptForTags(*repoFlag, tags)
	if err != nil {
		fmt.Printf("Error selecting tags: %v\n", err)
		os.Exit(1)
	}

	// Fetch commits between the selected tags
	commits, err := changelog.FetchCommits(*repoFlag, fromTag, toTag)
	if err != nil {
		fmt.Printf("Error fetching commits: %v\n", err)
		os.Exit(1)
	}

	// Prompt user to select commits to include
	selectedCommits, err := changelog.PromptForCommits(commits)
	if err != nil {
		fmt.Printf("Error selecting commits: %v\n", err)
		os.Exit(1)
	}

	// Generate the changelog content
	changelogContent := changelog.GenerateChangelog(selectedCommits, fromTag, toTag)

	// Render the changelog in Markdown format
	markdownChangelog, err := changelog.RenderMarkdown(changelogContent)
	if err != nil {
		fmt.Printf("Error rendering Markdown: %v\n", err)
		os.Exit(1)
	}

	// Render the changelog in HTML format
	htmlChangelog, err := changelog.RenderHTML(changelogContent)
	if err != nil {
		fmt.Printf("Error rendering HTML: %v\n", err)
		os.Exit(1)
	}

	// Write the Markdown changelog to CHANGELOG.md
	if err := os.WriteFile("CHANGELOG.md", []byte(markdownChangelog), 0644); err != nil {
		fmt.Printf("Error writing CHANGELOG.md: %v\n", err)
		os.Exit(1)
	}

	// Write the HTML changelog to public/release-notes.html
	if err := os.WriteFile("public/release-notes.html", []byte(htmlChangelog), 0644); err != nil {
		fmt.Printf("Error writing release-notes.html: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Changelog generated successfully!")
}
