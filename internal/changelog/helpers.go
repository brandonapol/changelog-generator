package changelog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetTagsForRepo prompts the user to enter the from and to tags for a repository
func GetTagsForRepo(repoPath string) (string, string, error) {
	fmt.Printf("Enter tags for %s:\n", repoPath)

	var fromTag, toTag string
	for {
		fmt.Print("From tag: ")
		if _, err := fmt.Scanln(&fromTag); err != nil {
			return "", "", fmt.Errorf("failed to read from tag: %w", err)
		}

		fmt.Print("To tag: ")
		if _, err := fmt.Scanln(&toTag); err != nil {
			return "", "", fmt.Errorf("failed to read to tag: %w", err)
		}

		// Validate tags
		if err := validateTag(repoPath, &fromTag); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		if err := validateTag(repoPath, &toTag); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		break
	}

	return fromTag, toTag, nil
}

func validateTag(repoPath string, tag *string) error {
	for {
		cmd := exec.Command("git", "rev-parse", *tag)
		cmd.Dir = repoPath
		if err := cmd.Run(); err != nil {
			fmt.Printf("Tag %s does not exist. Enter a new tag or list all tags? (n/l): ", *tag)
			var input string
			if _, err := fmt.Scanln(&input); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			if input == "l" {
				cmd := exec.Command("git", "tag")
				cmd.Dir = repoPath
				output, err := cmd.Output()
				if err != nil {
					return fmt.Errorf("failed to list tags: %w", err)
				}
				fmt.Printf("Available tags:\n%s\n", string(output))
			}

			fmt.Printf("Enter a new tag: ")
			if _, err := fmt.Scanln(tag); err != nil {
				return fmt.Errorf("failed to read tag: %w", err)
			}
		} else {
			return nil
		}
	}
}

// GenerateChangelog generates the changelog between two tags for a repository
func GenerateChangelog(repoPath, fromTag, toTag string) error {
	fmt.Printf("\nGenerating changelog for repository: %s\n", repoPath)

	// Get commit logs between the tags
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", fromTag, toTag), "--oneline", "--pretty=format:%s")
	cmd.Dir = repoPath
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	// TODO: Implement changelog generation logic
	// This is a placeholder for actual implementation

	return nil
}

// SelectCommits prompts the user to select commits to include in the changelog
func SelectCommits(repoPath, fromTag, toTag string) ([]string, error) {
	// Get commit logs between the tags
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", fromTag, toTag), "--oneline", "--pretty=format:%s")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	commits := strings.Split(string(output), "\n")
	// Filter out empty commits
	var filteredCommits []string
	for _, commit := range commits {
		if commit != "" {
			filteredCommits = append(filteredCommits, commit)
		}
	}
	commits = filteredCommits

	// Track selection state for each commit
	selected := make([]bool, len(commits))
	currentIndex := 0

	fmt.Println("Select commits to include in changelog:")
	fmt.Println("Enter: n (next), p (previous), s (select/deselect), d (done)")

	// Display initial view
	displayCommits(commits, selected, currentIndex)

	// Main selection loop
	for {
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)

		switch input {
		case "s", "S":
			// Toggle selection of current commit
			selected[currentIndex] = !selected[currentIndex]
			displayCommits(commits, selected, currentIndex)
		case "n", "N":
			// Move to next commit
			if currentIndex < len(commits)-1 {
				currentIndex++
			}
			displayCommits(commits, selected, currentIndex)
		case "p", "P":
			// Move to previous commit
			if currentIndex > 0 {
				currentIndex--
			}
			displayCommits(commits, selected, currentIndex)
		case "d", "D":
			// Done with selection
			var selectedCommits []string
			for i, isSelected := range selected {
				if isSelected {
					selectedCommits = append(selectedCommits, commits[i])
				}
			}
			return selectedCommits, nil
		case "q", "Q":
			// Quit
			return nil, fmt.Errorf("selection cancelled")
		default:
			fmt.Println("Invalid command. Use n (next), p (previous), s (select/deselect), d (done)")
		}
	}
}

// displayCommits displays the commit list with selection status
func displayCommits(commits []string, selected []bool, currentIndex int) {
	for i, commit := range commits {
		prefix := "○"
		if selected[i] {
			prefix = "●"
		}

		if i == currentIndex {
			fmt.Printf("-> %s %s\n", prefix, commit)
		} else {
			fmt.Printf("   %s %s\n", prefix, commit)
		}
	}
}

// UpdateChangelogFile updates the CHANGELOG.md file with the selected commits
func UpdateChangelogFile(changelogFile string, selectedCommits []string) error {
	// Read the existing CHANGELOG.md content
	content, err := os.ReadFile(changelogFile)
	if err != nil {
		return err
	}

	// Prepare the new changelog entry
	var newEntry strings.Builder
	newEntry.WriteString(fmt.Sprintf("## %s\n\n", time.Now().Format("2006-01-02")))
	for _, commit := range selectedCommits {
		newEntry.WriteString(fmt.Sprintf("- %s\n", commit))
	}
	newEntry.WriteString("\n")

	// Prepend the new entry to the existing content
	updatedContent := fmt.Sprintf("%s%s", newEntry.String(), string(content))

	// Write the updated content back to the file
	err = os.WriteFile(changelogFile, []byte(updatedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// UpdateReleaseNotesFile updates the release-notes.html file with the selected commits
func UpdateReleaseNotesFile(releaseNotesFile string, selectedCommits []string) error {
	// Read the existing release-notes.html content
	content, err := os.ReadFile(releaseNotesFile)
	if err != nil {
		return err
	}

	// Parse the HTML and find the <div class="notes-container"> element
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return err
	}

	notesContainer := doc.Find(".notes-container")

	// Create a new <div class="release-section"> element
	releaseSectionDoc, err := goquery.NewDocumentFromReader(strings.NewReader(`
		<div class="release-section">
			<div class="release-version">
				<span>Version X.Y.Z</span>
				<span class="release-date">Release Date</span>
			</div>
			<h3 class="change-category">Changes</h3>
			<ul class="change-list"></ul>
		</div>
	`))
	if err != nil {
		return err
	}
	releaseSection := releaseSectionDoc.Find(".release-section")

	// Update the version and release date
	releaseSection.Find(".release-version span").First().SetText(fmt.Sprintf("Version %s", "X.Y.Z")) // Replace with actual version
	releaseSection.Find(".release-date").SetText(time.Now().Format("January 02, 2006"))

	// Append the selected commits to the change list
	changeList := releaseSection.Find(".change-list")
	for _, commit := range selectedCommits {
		changeList.AppendHtml(fmt.Sprintf("<li>%s</li>", commit))
	}

	// Prepend the new release section to the notes container
	notesContainer.PrependSelection(releaseSection)

	// Write the updated HTML back to the file
	updatedContent, err := doc.Html()
	if err != nil {
		return err
	}

	err = os.WriteFile(releaseNotesFile, []byte(updatedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
