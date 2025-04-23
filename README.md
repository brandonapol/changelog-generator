# Changelog Generator

A command-line tool for generating changelogs and release notes from Git commits using conventional commit format.

## Prerequisites

- Go 1.20 or later
- Git

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/changelog-generator.git
   ```

2. Change to the project directory:

   ```
   cd changelog-generator
   ```

3. Build the executable:

   ```
   go build ./...
   ```

## Usage

To generate changelogs, run the following command:

```
./changelog-generator
```

The tool will interactively prompt you to:

1. Enter the names of the repositories you want to generate changelogs for (one per line, press Enter to finish)
2. For each repository:
   - Create a new tag (optional) based on the most recent tag
   - Select the 'from' and 'to' tags for the changelog
   - Select the commits to include in the changelog

The generated changelogs will be written to:
- `internal/changelog/output/CHANGELOG.md` (Markdown format)
- `internal/changelog/output/release-notes.html` (HTML format)

## Project Structure

- `cmd/changelog/main.go`: The main entry point for the CLI tool
- `internal/changelog/`
  - `generator.go`: Core changelog generation logic
  - `git.go`: Git-related functions for fetching commits and tags
  - `prompt.go`: Interactive prompt functions using the `promptui` package
  - `template.go`: HTML/Markdown template rendering functions
- `go.mod`, `go.sum`: Go module files

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Commit your changes
4. Push your branch to your forked repository
5. Open a pull request

Please ensure that your code follows the existing style and includes appropriate tests.

## License

This project is licensed under the [MIT License](LICENSE). 