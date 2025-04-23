package changelog

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var conventionalCommitTypes = "fix|feat|build|ci|docs|perf|refactor|revert|style|test"

// GenerateChangelog generates the changelog content from the given commits for multiple repositories
func GenerateChangelog(repositories []string, fromTag, toTag string) string {
	var changelog string
	for _, repo := range repositories {
		commits, err := FetchCommits(repo, fromTag, toTag)
		if err != nil {
			fmt.Printf("Error fetching commits for repository %s: %v\n", repo, err)
			continue
		}

		features, bugfixes, others := categorizeCommits(commits)

		// Format the categorized commits into changelog sections
		formattedFeatures := formatChangelogSection(features, "Features")
		formattedBugfixes := formatChangelogSection(bugfixes, "Bug Fixes")
		formattedOthers := formatChangelogSection(others, "Other Changes")

		// Prepare the changelog content for the repository
		repoChangelog := fmt.Sprintf("## Changelog for %s (%s)\n", repo, time.Now().Format("2006-01-02"))
		repoChangelog += fmt.Sprintf("%s → %s\n\n", fromTag, toTag)

		if formattedFeatures != "" {
			repoChangelog += fmt.Sprintf("### Features\n%s\n\n", formattedFeatures)
		}
		if formattedBugfixes != "" {
			repoChangelog += fmt.Sprintf("### Bug Fixes\n%s\n\n", formattedBugfixes)
		}
		if formattedOthers != "" {
			repoChangelog += fmt.Sprintf("### Other Changes\n%s\n\n", formattedOthers)
		}

		changelog += repoChangelog
	}
	return changelog
}

// categorizeCommits categorizes the commits into features, bugfixes and others
func categorizeCommits(commits []string) (features, bugfixes, others []string) {
	// Regular expression to check for valid conventional commit formats
	re := regexp.MustCompile(fmt.Sprintf(`^(%s)\([A-Za-z0-9-]*\):|^(%s):`, conventionalCommitTypes, conventionalCommitTypes))

	for _, commit := range commits {
		if !re.MatchString(commit) {
			fmt.Printf("Skipping commit: %s\n", commit)
			continue
		}

		// Strip the scope in the conventional commit
		commitMessage := regexp.MustCompile(`\([A-Za-z0-9-]*\):`).ReplaceAllString(commit, ":")
		// Strip leading and trailing space
		commitMessage = strings.TrimSpace(commitMessage)
		fmt.Printf("Adding commit: %s\n", commitMessage)

		// Categorize based on commit type
		switch {
		case strings.HasPrefix(commitMessage, "fix:"):
			bugfixes = append(bugfixes, commitMessage)
		case strings.HasPrefix(commitMessage, "feat:"):
			features = append(features, commitMessage)
		default:
			others = append(others, commitMessage)
		}
	}

	return
}

// formatChangelogSection formats the commits into a changelog section
func formatChangelogSection(commits []string, sectionHeading string) string {
	if len(commits) == 0 {
		return ""
	}

	formattedSection := fmt.Sprintf("### %s\n", sectionHeading)
	for _, commit := range commits {
		formattedSection += fmt.Sprintf("- %s\n", strings.TrimPrefix(commit, strings.Split(commit, ":")[0]+": "))
	}

	return formattedSection
}
