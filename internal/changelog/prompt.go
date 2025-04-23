package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptForTags prompts the user to select the from and to tags
func PromptForTags(repoPath string, tags []string) (string, string, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Select the 'from' tag for %s", repoPath),
		Items: tags,
	}
	_, fromTag, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to select 'from' tag: %v", err)
	}

	prompt = promptui.Select{
		Label: fmt.Sprintf("Select the 'to' tag for %s", repoPath),
		Items: tags,
	}
	_, toTag, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to select 'to' tag: %v", err)
	}

	return fromTag, toTag, nil
}

// PromptForCommits prompts the user to select the commits to include in the changelog
func PromptForCommits(commits []string) ([]string, error) {
	var selectedCommits []string
	for _, commit := range commits {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Include commit '%s' in changelog?", commit),
			IsConfirm: true,
		}
		result, err := prompt.Run()
		if err != nil && err != promptui.ErrAbort {
			return nil, fmt.Errorf("failed to prompt for commit: %v", err)
		}
		if strings.ToLower(result) == "y" {
			selectedCommits = append(selectedCommits, commit)
		}
	}
	return selectedCommits, nil
}

// PromptForNewTag prompts the user if they want to create a new tag based on the most recent tag
func PromptForNewTag(mostRecentTag string) (bool, string, error) {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Last tag was %s, make a new tag?", mostRecentTag),
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false, "", fmt.Errorf("failed to prompt for new tag: %v", err)
	}
	createNewTag := strings.ToLower(result) == "y"

	if createNewTag {
		prompt := promptui.Prompt{
			Label: "Enter new tag name",
		}
		newTag, err := prompt.Run()
		if err != nil {
			return false, "", fmt.Errorf("failed to prompt for new tag name: %v", err)
		}
		return true, newTag, nil
	}

	return false, "", nil
}

// PromptForRepositories prompts the user to enter the list of repositories to generate the changelog for
func PromptForRepositories() ([]string, error) {
	var repositories []string
	fmt.Println("Enter the paths to your Git repositories.")
	fmt.Println("You can use relative paths:")
	fmt.Println("  - '.' for the current directory")
	fmt.Println("  - '../repo-name' for a sibling directory")
	fmt.Println("  - Or any other relative or absolute path")

	for {
		prompt := promptui.Prompt{
			Label: "Enter a repository path (or press Enter to finish)",
		}
		result, err := prompt.Run()
		if err != nil {
			return nil, fmt.Errorf("failed to prompt for repository: %v", err)
		}
		if result == "" {
			break
		}

		// Convert relative path to absolute path
		absPath, err := filepath.Abs(result)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve path %s: %v", result, err)
		}

		// Verify that the path exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Printf("Warning: Path %s does not exist. Please enter a valid path.\n", absPath)
			continue
		}

		repositories = append(repositories, absPath)
		fmt.Printf("Added repository: %s\n", absPath)
	}

	return repositories, nil
}
