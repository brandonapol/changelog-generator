package changelog

import (
	"fmt"
	"regexp"
	"strings"
)

var conventionalCommitTypes = "fix|feat|build|ci|docs|perf|refactor|revert|style|test"

// GenerateChangelog generates the changelog content from the given commits for multiple repositories
func GenerateChangelog(repoCommits map[string][]string) string {
	var features, bugfixes, others []string
	for _, commits := range repoCommits {
		repoFeatures, repoBugfixes, repoOthers := categorizeCommits(commits)
		features = append(features, repoFeatures...)
		bugfixes = append(bugfixes, repoBugfixes...)
		others = append(others, repoOthers...)
	}

	// Format the categorized commits into changelog sections
	formattedFeatures := formatChangelogSection(features, "Features")
	formattedBugfixes := formatChangelogSection(bugfixes, "Bug Fixes")
	formattedOthers := formatChangelogSection(others, "Other Changes")

	// Prepare the combined changelog content
	changelog := "## Combined Changelog\n\n"
	if formattedFeatures != "" {
		changelog += fmt.Sprintf("### Features\n%s\n\n", formattedFeatures)
	}
	if formattedBugfixes != "" {
		changelog += fmt.Sprintf("### Bug Fixes\n%s\n\n", formattedBugfixes)
	}
	if formattedOthers != "" {
		changelog += fmt.Sprintf("### Other Changes\n%s\n\n", formattedOthers)
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
