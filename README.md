# Changelog Generator

A command-line tool for generating changelogs and release notes from Git commits using conventional commit format.

## Prerequisites

- Go 1.20 or later
- Git
- Make (optional, for using the Makefile)

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
   make build
   ```

   Or with Go directly:

   ```
   go build -o bin/changelog-generator ./cmd/changelog
   ```

## Usage

To generate changelogs, run the following command:

```
make run
```

Or with the binary directly:

```
./bin/changelog-generator
```

The tool will interactively prompt you to:

1. Enter the paths to the repositories you want to generate changelogs for (one per line, press Enter to finish)
   - You can use relative paths:
     - `.` for the current directory
     - `../repo-name` for a sibling directory
     - Or any other relative or absolute path
2. For each repository:
   - Create a new tag (optional) based on the most recent tag
   - Select the 'from' and 'to' tags for the changelog
   - Select the commits to include in the changelog

The generated changelogs will be written to:
- `internal/changelog/output/CHANGELOG.md` (Markdown format)
- `internal/changelog/output/release-notes.html` (HTML format)

## Makefile Commands

The project includes a Makefile with the following commands:

- `make build` - Build the binary
- `make clean` - Remove build artifacts
- `make deps` - Install dependencies
- `make lint` - Run linter
- `make test` - Run tests
- `make run` - Run the generator
- `make help` - Show help message

## Project Structure

- `cmd/changelog/main.go`: The main entry point for the CLI tool
- `internal/changelog/`
  - `generator.go`: Core changelog generation logic
  - `git.go`: Git-related functions for fetching commits and tags
  - `prompt.go`: Interactive prompt functions using the `promptui` package
  - `template.go`: HTML/Markdown template rendering functions
- `go.mod`, `go.sum`: Go module files
- `Makefile`: Build automation
- `.github/workflows/go.yml`: GitHub Actions CI workflow

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Commit your changes
4. Push your branch to your forked repository
5. Open a pull request

Please ensure that your code follows the existing style and passes the linter checks.

## License

This project is licensed under the [MIT License](LICENSE). 