package changelog

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

// RenderMarkdown renders the changelog in Markdown format and writes it to the provided file path
func RenderMarkdown(changelog, appVersion, releaseDate, changelogFile string) error {
	if changelogFile == "" {
		return nil
	}

	// Read existing CHANGELOG.md content
	existingContent, err := os.ReadFile(changelogFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %v", changelogFile, err)
	}

	// Prepare the new changelog content
	newContent := fmt.Sprintf("## Changelog (%s)\nApp Version: %s\n\n%s\n", releaseDate, appVersion, changelog)

	// Prepend the new content to the existing content
	finalContent := newContent + string(existingContent)

	// Write the updated content to CHANGELOG.md
	if err := os.WriteFile(changelogFile, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", changelogFile, err)
	}

	return nil
}

// RenderHTML renders the changelog in HTML format and writes it to the provided file path
func RenderHTML(changelog string, appVersion, releaseDate string, features, bugfixes, others []string, releaseNotesFile string) error {
	if releaseNotesFile == "" {
		return nil
	}

	// Parse the existing release-notes.html file
	existingContent, err := os.ReadFile(releaseNotesFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %v", releaseNotesFile, err)
	}

	// Extract the existing release sections
	existingSections := ""
	if len(existingContent) > 0 {
		startIndex := bytes.Index(existingContent, []byte("<div class=\"notes-container\">"))
		endIndex := bytes.LastIndex(existingContent, []byte("</div>"))

		// Make sure indices are valid and in the right order
		if startIndex != -1 && endIndex != -1 && startIndex+30 < endIndex && endIndex < len(existingContent) {
			existingSections = string(existingContent[startIndex+30 : endIndex])
		} else {
			// If we can't find the proper structure, just initialize with empty sections
			fmt.Printf("Warning: Could not parse existing sections in %s, starting with empty content\n", releaseNotesFile)
		}
	}

	// Generate the new release section
	newSection := fmt.Sprintf(`
		<div class="release-section">
			<div class="release-version">
				<span>Version %s</span>
				<span class="release-date">%s</span>
			</div>
			
			<h3 class="change-category">Features</h3>
			<ul class="change-list">
				%s
			</ul>
			
			<h3 class="change-category">Bug Fixes</h3>
			<ul class="change-list">
				%s
			</ul>
			
			<h3 class="change-category">Other Changes</h3>
			<ul class="change-list">
				%s
			</ul>
		</div>
	`, appVersion, releaseDate, formatChangeList(features), formatChangeList(bugfixes), formatChangeList(others))

	// Combine the existing sections with the new section
	combinedSections := newSection + existingSections

	// Parse the HTML template
	tmpl, err := template.ParseFS(templateFS, "templates/release-notes.html")
	if err != nil {
		return err
	}

	// Render the template with the updated release sections
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"AppName":     "MyApp",
		"AppVersion":  appVersion,
		"ReleaseDate": releaseDate,
		"Changelog":   template.HTML(combinedSections),
	}); err != nil {
		return err
	}

	// Write the updated release notes to file
	if err := os.WriteFile(releaseNotesFile, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", releaseNotesFile, err)
	}

	return nil
}

// formatChangeList formats a list of changes as HTML list items
func formatChangeList(changes []string) string {
	if len(changes) == 0 {
		return "<li>No changes</li>"
	}

	var sb strings.Builder
	for _, change := range changes {
		sb.WriteString("<li>")
		sb.WriteString(change)
		sb.WriteString("</li>\n")
	}
	return sb.String()
}
