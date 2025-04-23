package changelog

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// FetchCommits fetches the commits between two tags
func FetchCommits(repoPath, fromTag, toTag string) ([]string, error) {
	// Check if the repository path exists and is a git repository
	if err := validateRepoPath(repoPath); err != nil {
		return nil, err
	}

	// Navigate to the repository
	cmd := exec.Command("git", "-C", repoPath, "log", "--oneline", "--pretty=format:%s", fmt.Sprintf("%s..%s", fromTag, toTag))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch commits: %v", err)
	}

	commits := strings.Split(string(output), "\n")
	// Remove empty commits
	var filteredCommits []string
	for _, commit := range commits {
		if commit != "" {
			filteredCommits = append(filteredCommits, commit)
		}
	}

	return filteredCommits, nil
}

// FetchTags fetches all the tags in the repository
func FetchTags(repoPath string) ([]string, error) {
	// Check if the repository path exists and is a git repository
	if err := validateRepoPath(repoPath); err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "-C", repoPath, "tag")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %v", err)
	}

	tags := strings.Split(string(output), "\n")
	// Remove empty tags
	var filteredTags []string
	for _, tag := range tags {
		if tag != "" {
			filteredTags = append(filteredTags, tag)
		}
	}

	// Check if we have any tags
	if len(filteredTags) == 0 {
		return nil, fmt.Errorf("no tags found in repository %s", repoPath)
	}

	return filteredTags, nil
}

// ValidateTag checks if a tag exists in the repository
func ValidateTag(repoPath, tag string) error {
	// Check if the repository path exists and is a git repository
	if err := validateRepoPath(repoPath); err != nil {
		return err
	}

	cmd := exec.Command("git", "-C", repoPath, "rev-parse", tag)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tag %s does not exist in %s", tag, repoPath)
	}
	return nil
}

// CreateTag creates a new tag in the repository
func CreateTag(repoPath, tag string) error {
	// Check if the repository path exists and is a git repository
	if err := validateRepoPath(repoPath); err != nil {
		return err
	}

	cmd := exec.Command("git", "-C", repoPath, "tag", tag)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag %s: %v", tag, err)
	}
	return nil
}

// validateRepoPath checks if the provided path exists and is a git repository
func validateRepoPath(repoPath string) error {
	// Check if path exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path %s does not exist", repoPath)
	}

	// Check if path is a git repository
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("path %s is not a valid git repository - make sure you're inside a git repository or specify a valid path", repoPath)
	}

	return nil
}
