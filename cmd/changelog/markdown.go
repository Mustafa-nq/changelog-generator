package main

import (
	"fmt"
	"os"
	"time"

	"changelog-generator/internal/lib"
)

// GenerateMarkdown creates a formatted markdown changelog
func GenerateMarkdown(commits []*lib.Commit, projectName, version string) string {
	// Start with header
	md := fmt.Sprintf("# Changelog - %s\n\n", projectName)
	md += fmt.Sprintf("## Version %s\n", version)
	md += fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("January 2, 2006"))

	// Group commits by category
	groups := lib.GroupCommitsByCategory(commits)

	// Define order of categories
	categories := []lib.CommitCategory{
		lib.CategoryBreaking,
		lib.CategoryFeature,
		lib.CategoryFix,
		lib.CategoryPerformance,
		lib.CategoryRefactor,
		lib.CategoryDocs,
		lib.CategoryTest,
		lib.CategoryChore,
		lib.CategoryOther,
	}

	// Add each category
	for _, category := range categories {
		categoryCommits, exists := groups[category]
		if !exists || len(categoryCommits) == 0 {
			continue // Skip empty categories
		}

		// Category header
		md += fmt.Sprintf("### %s\n\n", category)

		// List commits
		for _, commit := range categoryCommits {
			// Clean up the commit message (remove prefixes)
			message := cleanCommitMessage(commit.Message)
			md += fmt.Sprintf("- %s ([%s])\n", message, commit.Hash)
		}

		md += "\n"
	}

	// Add footer
	md += "---\n"
	md += fmt.Sprintf("*Total commits: %d*\n", len(commits))

	return md
}

// cleanCommitMessage removes conventional commit prefixes
func cleanCommitMessage(msg string) string {
	// Remove trailing newlines
	cleaned := ""
	for _, c := range msg {
		if c == '\n' {
			break
		}
		cleaned += string(c)
	}

	// Remove conventional commit prefix (feat:, fix:, etc.)
	prefixes := []string{"feat:", "fix:", "docs:", "chore:", "test:", "refactor:", "perf:"}
	for _, prefix := range prefixes {
		if len(cleaned) > len(prefix) && cleaned[:len(prefix)] == prefix {
			// Remove prefix and trim space
			cleaned = cleaned[len(prefix):]
			// Remove leading space
			if len(cleaned) > 0 && cleaned[0] == ' ' {
				cleaned = cleaned[1:]
			}
			break
		}
	}

	// Remove scope: feat(auth): -> just remove the (auth) part
	// Find and remove text between ( and ):
	result := ""
	skipUntil := rune(0)
	for _, c := range cleaned {
		if c == '(' && skipUntil == 0 {
			skipUntil = ')'
			continue
		}
		if skipUntil != 0 {
			if c == skipUntil {
				skipUntil = 0
				// Skip the colon and space after )
				continue
			}
			continue
		}
		result += string(c)
	}

	// Remove leading ': '
	if len(result) > 2 && result[:2] == ": " {
		result = result[2:]
	}

	// Capitalize first letter
	if len(result) > 0 {
		first := result[0]
		if first >= 'a' && first <= 'z' {
			result = string(first-32) + result[1:]
		}
	}

	return result
}

// SaveMarkdown saves the markdown content to a file
func SaveMarkdown(content, filename string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
