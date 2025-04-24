package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brandonapol/changelog-generator/internal/changelog"
	"github.com/fatih/color"
)

func main() {
	color.Cyan("Enter the name of your project:")
	var projectName string
	if _, err := fmt.Scanln(&projectName); err != nil {
		color.Red("Error reading project name: %v", err)
		os.Exit(1)
	}

	var repoPaths []string
	for {
		color.Cyan("Enter a repository path (or press Enter to finish):")
		var repoPath string
		if _, err := fmt.Scanln(&repoPath); err != nil {
			// EOF or unexpected newline might indicate empty input
			if err.Error() == "unexpected newline" || err.Error() == "EOF" {
				break
			}
			color.Red("Error reading repository path: %v", err)
			continue
		}
		if repoPath == "" {
			break
		}
		absRepoPath, err := filepath.Abs(repoPath)
		if err != nil {
			color.Red("Error: Invalid repository path %s\n", repoPath)
			continue
		}
		repoPaths = append(repoPaths, absRepoPath)
		color.Green("Added repository: %s\n", absRepoPath)
	}

	if len(repoPaths) == 0 {
		color.Red("No repositories provided. Exiting.")
		os.Exit(1)
	}

	// Generate changelog for each repository
	for _, repoPath := range repoPaths {
		// Get tags for the repository
		color.Cyan("Enter tags for %s:", repoPath)
		fromTag, toTag, err := changelog.GetTagsForRepo(repoPath)
		if err != nil {
			color.Red("Error getting tags for repository %s: %v\n", repoPath, err)
			continue
		}

		// Generate changelog between the tags
		color.Cyan("Generating changelog for repository: %s\n", repoPath)
		err = changelog.GenerateChangelog(repoPath, fromTag, toTag)
		if err != nil {
			color.Red("Error generating changelog for repository %s: %v\n", repoPath, err)
			continue
		}

		// Find CHANGELOG.md and release-notes.html files recursively
		changelogFile, releaseNotesFile, err := findChangelogFiles(repoPath)
		if err != nil {
			color.Red("Error finding changelog files for repository %s: %v\n", repoPath, err)
			continue
		}

		// Prompt user to select commits
		color.Cyan("Select commits to include in changelog:")
		selectedCommits, err := changelog.SelectCommits(repoPath, fromTag, toTag)
		if err != nil {
			color.Red("Error selecting commits for repository %s: %v\n", repoPath, err)
			continue
		}

		// Update CHANGELOG.md
		err = changelog.UpdateChangelogFile(changelogFile, selectedCommits)
		if err != nil {
			color.Red("Error updating CHANGELOG.md for repository %s: %v\n", repoPath, err)
			continue
		}

		// Update release-notes.html
		err = changelog.UpdateReleaseNotesFile(releaseNotesFile, selectedCommits)
		if err != nil {
			color.Red("Error updating release-notes.html for repository %s: %v\n", repoPath, err)
			continue
		}

		color.Green("Changelog generated successfully for repository: %s\n", repoPath)
	}
}

func findChangelogFiles(repoPath string) (string, string, error) {
	var changelogFile, releaseNotesFile string

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && (info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".")) {
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
		return "", "", err
	}

	if changelogFile == "" {
		return "", "", fmt.Errorf("CHANGELOG.md not found in repository %s", repoPath)
	}

	if releaseNotesFile == "" {
		return "", "", fmt.Errorf("release-notes.html not found in repository %s", repoPath)
	}

	return changelogFile, releaseNotesFile, nil
}
