package lib

import (
	"fmt"

	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Represent a git commit
type Commit struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
}

// CommitCategory represents the type of change a commit introduces (feature, bugfix, etc.)
type CommitCategory string

const (
	CategoryFeature     CommitCategory = " Features"
	CategoryFix         CommitCategory = " Bug Fixes"
	CategoryBreaking    CommitCategory = "  Breaking Changes"
	CategoryDocs        CommitCategory = " Documentation"
	CategoryPerformance CommitCategory = " Performance"
	CategoryRefactor    CommitCategory = "  Refactoring"
	CategoryTest        CommitCategory = " Tests"
	CategoryChore       CommitCategory = " Chores"
	CategoryOther       CommitCategory = " Other"
)

// CategorizeCommit categorizes a commit based on its message
func CategorizeCommit(commit *Commit) CommitCategory {
	msg := commit.Message

	// Check for breaking change indicator
	if containsAny(msg, []string{"BREAKING CHANGE", "BREAKING:", "!"}) {
		return CategoryBreaking
	}

	// Check conventional commit prefixes
	if hasPrefix(msg, "feat") {
		return CategoryFeature
	}
	if hasPrefix(msg, "fix") {
		return CategoryFix
	}
	if hasPrefix(msg, "docs") {
		return CategoryDocs
	}
	if hasPrefix(msg, "perf") {
		return CategoryPerformance
	}
	if hasPrefix(msg, "refactor") {
		return CategoryRefactor
	}
	if hasPrefix(msg, "test") {
		return CategoryTest
	}
	if hasPrefix(msg, "chore") {
		return CategoryChore
	}

	// Check for keywords in message
	if containsAny(msg, []string{"add", "implement", "create", "new"}) {
		return CategoryFeature
	}
	if containsAny(msg, []string{"fix", "bug", "issue", "resolve"}) {
		return CategoryFix
	}
	if containsAny(msg, []string{"update", "improve"}) {
		return CategoryRefactor
	}

	return CategoryOther
}

// hasPrefix checks if message starts with a conventional commit prefix
func hasPrefix(msg, prefix string) bool {
	// Check patterns like "feat:", "feat(scope):", "feat!"
	patterns := []string{
		prefix + ":",
		prefix + "(",
		prefix + "!",
	}

	for _, pattern := range patterns {
		if len(msg) >= len(pattern) && msg[:len(pattern)] == pattern {
			return true
		}
	}
	return false
}

// containsAny checks if message contains any of the keywords
func containsAny(msg string, keywords []string) bool {
	msgLower := toLower(msg)
	for _, keyword := range keywords {
		if contains(msgLower, toLower(keyword)) {
			return true
		}
	}
	return false
}

// Simple string helpers
func toLower(s string) string {
	result := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			result += string(c + 32)
		} else {
			result += string(c)
		}
	}
	return result
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GroupCommitsByCategory groups commits by their category
func GroupCommitsByCategory(commits []*Commit) map[CommitCategory][]*Commit {
	groups := make(map[CommitCategory][]*Commit)

	for _, commit := range commits {
		category := CategorizeCommit(commit)
		groups[category] = append(groups[category], commit)
	}

	return groups
}

// PrintGroupedCommits displays commits grouped by category
func PrintGroupedCommits(groups map[CommitCategory][]*Commit) {
	// Define order of categories
	categories := []CommitCategory{
		CategoryBreaking,
		CategoryFeature,
		CategoryFix,
		CategoryPerformance,
		CategoryRefactor,
		CategoryDocs,
		CategoryTest,
		CategoryChore,
		CategoryOther,
	}

	fmt.Println("Categorized Commits:")
	fmt.Println()

	for _, category := range categories {
		commits, exists := groups[category]
		if !exists || len(commits) == 0 {
			continue // Skip empty categories
		}

		fmt.Printf("%s (%d)\n", category, len(commits))
		fmt.Println("─────────────────────────────────────")

		for _, commit := range commits {
			fmt.Printf("  [%s] %s\n", commit.Hash, commit.Message)
		}
		fmt.Println()
	}
}

// OpenRepository opens a git repository at the given path
func OpenRepository(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}
	return repo, nil
}

// GetRecentCommits gets the last N commits from the repository
func GetRecentCommits(repo *git.Repository, count int) ([]*Commit, error) {
	// Get HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Get commit history starting from HEAD
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Collect commits
	var commits []*Commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		// Stop if we have enough commits
		if len(commits) >= count {
			return fmt.Errorf("done") // Use error to break iteration
		}

		// Add commit to our list
		commits = append(commits, &Commit{
			Hash:    c.Hash.String()[:7], // Short hash (first 7 chars)
			Author:  c.Author.Name,
			Date:    c.Author.When,
			Message: c.Message,
		})

		return nil
	})

	// Ignore the "done" error we used to break the loop
	if err != nil && err.Error() != "done" {
		return nil, fmt.Errorf("error iterating commits: %w", err)
	}

	return commits, nil
}

func PrintCommits(commits []*Commit) {
	fmt.Println("Recent Commits: ")
	fmt.Println()

	for i, commit := range commits {
		fmt.Printf("%d. [%s] %s\n", i+1, commit.Hash, commit.Author)
		fmt.Printf("   %s\n", commit.Date.Format("2006-01-02 15:04:05"))
		fmt.Printf("   %s\n", commit.Message)
		fmt.Println()
	}
}
