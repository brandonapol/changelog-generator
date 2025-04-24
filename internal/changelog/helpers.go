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

const (
	conventionalCommitTypes = "fix|feat|build|ci|docs|perf|refactor|revert|style|test"
)

// GetTagsForRepo prompts the user to enter the from and to tags for a repository
func GetTagsForRepo(repoPath string) (string, string, error) {
	fmt.Printf("Enter tags for %s:\n", repoPath)

	var fromTag, toTag string
	fmt.Print("From tag: ")
	fmt.Scanln(&fromTag)
	fmt.Print("To tag: ")
	fmt.Scanln(&toTag)

	// Validate tags exist
	cmd := exec.Command("git", "rev-parse", fromTag)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("tag %s does not exist in %s", fromTag, repoPath)
	}

	cmd = exec.Command("git", "rev-parse", toTag)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("tag %s does not exist in %s", toTag, repoPath)
	}

	return fromTag, toTag, nil
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

	fmt.Println("Select commits to include in changelog (Space to select, Enter when done):")
	var selectedCommits []string
	for _, commit := range commits {
		if commit == "" {
			continue
		}

		fmt.Printf("○ %s", commit)
		var input string
		fmt.Scanln(&input)
		if input == " " {
			selectedCommits = append(selectedCommits, commit)
			fmt.Printf("● %s\n", commit)
		} else {
			fmt.Println()
		}
	}

	return selectedCommits, nil
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
