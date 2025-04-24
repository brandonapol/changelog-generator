// Package changelog provides functionality for generating changelogs from git repositories
package changelog

import (
	"github.com/brandonapol/changelog-generator/internal/changelog"
)

// GetTagsForRepo prompts the user to enter tags for a given repository
func GetTagsForRepo(repoPath string) (string, string, error) {
	return changelog.GetTagsForRepo(repoPath)
}

// GenerateChangelog generates a changelog for a repository between two tags
func GenerateChangelog(repoPath, fromTag, toTag string) error {
	return changelog.GenerateChangelog(repoPath, fromTag, toTag)
}

// SelectCommits allows the user to select commits to include in the changelog
func SelectCommits(repoPath, fromTag, toTag string) ([]string, error) {
	return changelog.SelectCommits(repoPath, fromTag, toTag)
}

// UpdateChangelogFile updates the CHANGELOG.md file with the selected commits
func UpdateChangelogFile(changelogFile string, commits []string) error {
	return changelog.UpdateChangelogFile(changelogFile, commits)
}

// UpdateReleaseNotesFile updates the release-notes.html file with the selected commits
func UpdateReleaseNotesFile(releaseNotesFile string, commits []string) error {
	return changelog.UpdateReleaseNotesFile(releaseNotesFile, commits)
}
