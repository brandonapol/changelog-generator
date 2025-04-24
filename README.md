# Changelog Generator

A command-line tool for generating changelogs and release notes from Git repositories using conventional commit format.

## Features

- Generate changelogs between Git tags
- Interactive commit selection
- Automatic updating of CHANGELOG.md and release-notes.html files
- Support for multiple repositories
- Conventional commit detection

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/brandonapol/changelog-generator.git
cd changelog-generator

# Build the project
go build -o bin/changelog-generator ./cmd/changelog

# Add to your PATH (optional)
cp bin/changelog-generator /usr/local/bin/
```

## Usage

Run the tool and follow the interactive prompts:

```bash
changelog-generator
```

The tool will:

1. Ask for your project name
2. Prompt for repository paths (you can add multiple)
3. For each repository:
   - Ask for the from/to tags for the changelog
   - Generate a changelog between those tags
   - Allow you to select specific commits to include
   - Update CHANGELOG.md and release-notes.html files

## Project Structure

```
changelog-generator/
├── cmd/
│   └── changelog/      # Command-line application
│       └── main.go     # Entry point
├── internal/
│   └── changelog/      # Internal implementation
│       ├── helpers.go  # Core functionality
│       └── templates/  # HTML/Markdown templates
├── changelog.go        # Public API
└── go.mod              # Go module definition
```

## Contributing

Contributions are welcome! Here's how you can help:

1. **Code Contributions**
   - Fork the repository
   - Create a feature branch (`git checkout -b feature/amazing-feature`)
   - Commit your changes (`git commit -m 'Add amazing feature'`)
   - Push to the branch (`git push origin feature/amazing-feature`)
   - Open a Pull Request

2. **Bug Reports**
   - Use the issue tracker to report bugs
   - Include detailed steps to reproduce the issue
   - Mention your environment (OS, Go version)

3. **Feature Requests**
   - Use the issue tracker to suggest features
   - Describe the feature in detail and its use cases

### Development Guidelines

- Follow idiomatic Go practices
- Prefer functional programming patterns when appropriate
- Maintain backward compatibility for public APIs
- Write tests for new functionality
- Document public functions and types

## License

This project is licensed under the MIT License - see the LICENSE file for details.

