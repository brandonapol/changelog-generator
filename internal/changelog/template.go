package changelog

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"time"
)

//go:embed templates/*
var templateFS embed.FS

// RenderMarkdown renders the changelog in Markdown format and writes it to internal/changelog/output/CHANGELOG.md
func RenderMarkdown(changelog string) error {
	tmpl, err := template.New("markdown").Parse(`# Changelog

{{ .Changelog }}
`)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"Changelog": changelog,
	}); err != nil {
		return err
	}

	if err := os.MkdirAll("internal/changelog/output", os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile("internal/changelog/output/CHANGELOG.md", buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

// RenderHTML renders the changelog in HTML format and writes it to internal/changelog/output/release-notes.html
func RenderHTML(changelog string) error {
	tmpl, err := template.ParseFS(templateFS, "templates/release-notes.html")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"Changelog":   template.HTML(changelog),
		"AppName":     "MyApp",
		"AppVersion":  "3.0.3",
		"ReleaseDate": time.Now().Format("January 2, 2006"),
	}); err != nil {
		return err
	}

	if err := os.MkdirAll("internal/changelog/output", os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile("internal/changelog/output/release-notes.html", buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
