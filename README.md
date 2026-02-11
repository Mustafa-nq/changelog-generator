# Changelog Generator

An AI-powered changelog generator that transforms your Git commit history into beautiful, professional release notes.


## Features

-  **AI-Powered** - Uses Claude AI to improve commit messages for clarity
-  **Smart Categorization** - Automatically categorizes commits (features, fixes, breaking changes, etc.)
-  **Beautiful Output** - Generates clean, formatted Markdown changelogs
-  **Conventional Commits** - Supports conventional commit format
-  **Configurable** - YAML-based configuration
-  **Fast** - Built with Go for performance

##  Example Output
```markdown
# Changelog - My Project

## Version 1.0.0
**Generated:** February 6, 2024

### ✨ Features

- Add user authentication system ([abc123d])
- Implement password reset functionality ([def456a])

### 🐛 Bug Fixes

- Fix memory leak in background workers ([789ghij])
- Resolve crash on invalid input ([klm012n])

### 📝 Documentation

- Update installation guide ([opq345r])

---
*Total commits: 5*
```

##  Why Use This?

- **Save Time** - Generate changelogs in seconds, not hours
- **Consistency** - Always get the same professional format
- **User-Friendly** - AI rewrites technical commits for end users
- **Automatic** - Just run one command

##  Installation

### Prerequisites

- Go 1.21 or higher
- Git
- (Optional) Claude API key for AI features

### Install
```bash
# Clone the repository
git clone https://github.com/yourusername/changelog-generator.git
cd changelog-generator

# Download dependencies
go mod download

# Build
go build -o changelog ./cmd/changelog

# (Optional) Install globally
sudo mv changelog /usr/local/bin/
```

##  Quick Start

### 1. Initialize Configuration
```bash
# Navigate to your Git repository
cd /path/to/your/project

# Create config file
changelog init
```

This creates a `.changelogrc.yaml` file with default settings.

### 2. Generate Changelog
```bash
# Generate from last 10 commits
changelog generate

# Generate from last 20 commits
changelog generate --count 20

# Use AI to improve messages (requires API key)
export CLAUDE_API_KEY="your-api-key-here"
changelog generate --ai
```

### 3. View Output

Your changelog is saved to `CHANGELOG.md` (or the filename specified in config).
```bash
cat CHANGELOG.md
```

## Usage

### Commands
```bash
changelog init              # Initialize configuration
changelog generate          # Generate changelog
changelog generate --ai     # Generate with AI improvements
changelog show             # Show current configuration
changelog --help           # Show all commands
changelog --version        # Show version
```

### Flags
```bash
--count N         # Number of commits to include (default: 10)
--output FILE     # Output filename (default: from config)
--since REF       # Starting point (default: HEAD~10)
--to REF          # Ending point (default: HEAD)
--ai              # Use AI to improve commit messages
```

### Examples
```bash
# Generate from last release
changelog generate --since v1.0.0 --to HEAD

# Generate with AI enhancement
changelog generate --count 20 --ai

# Custom output file
changelog generate --output RELEASE_NOTES.md

# Specific range
changelog generate --since abc123 --to def456
```

##  Configuration

Edit `.changelogrc.yaml` to customize:
```yaml
# Project information
project:
  name: "My Project"
  version: "1.0.0"

# Git settings
git:
  repository_path: "."
  default_branch: "main"

# Output settings
output:
  format: "markdown"
  filename: "CHANGELOG.md"

# Commit categories
categories:
  - breaking
  - features
  - fixes
  - documentation
```

## 🤖 AI Features

The tool can use Claude AI to improve commit messages:

### Setup

1. Get an API key from [Anthropic](https://console.anthropic.com/)
2. Set environment variable:
```bash
   export CLAUDE_API_KEY="sk-ant-your-key-here"
```
3. Use the `--ai` flag:
```bash
   changelog generate --ai
```

##  How It Works

1. **Scans Git History** - Reads commits from your repository
2. **Categorizes** - Groups commits by type (feature, fix, etc.)
3. **Improves** (optional) - Uses AI to rewrite messages clearly
4. **Formats** - Generates beautiful Markdown
5. **Saves** - Writes to CHANGELOG.md


##  Tech Stack

- **Language:** Go 1.21+
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
- **Git Library:** [go-git](https://github.com/go-git/go-git)
- **AI:** LLM API
- **Config:** YAML

##  Development

### Build
```bash
go build -o changelog ./cmd/changelog
```
