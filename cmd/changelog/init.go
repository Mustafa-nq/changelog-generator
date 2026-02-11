package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize changelog configuration",
	Long: `Create a configuration file for the changelog generator.

This will set up everything you need to start generating changelogs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing changelog configuration...")
		fmt.Println()

		// Create the config file
		err := createConfigFile()
		if err != nil {
			fmt.Printf("Error creating config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Created .changelogrc.yaml")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Edit .changelogrc.yaml to set your preferences")
		fmt.Println("  2. Run 'changelog generate' to create your first changelog")
	},
}

// createConfigFile creates the default configuration file
func createConfigFile() error {
	// Check if file already exists
	if _, err := os.Stat(".changelogrc.yaml"); err == nil {
		// File exists
		fmt.Println(".changelogrc.yaml already exists!")
		fmt.Print("Overwrite? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}
	}

	// Default configuration content
	configContent := `# Changelog Generator Configuration

# Your project information
project:
  name: "My Project"
  version: "1.0.0"

# Git repository settings
git:
  repository_path: "."
  default_branch: "main"

# Output settings
output:
  format: "markdown"
  filename: "CHANGELOG.md"

# AI settings (for future use)
ai:
  enabled: false
  provider: "claude"

# Categories for changes
categories:
  - breaking
  - features
  - fixes
  - documentation
`

	// Write the file
	err := os.WriteFile(".changelogrc.yaml", []byte(configContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func init() {
	// Add init command to root command
	rootCmd.AddCommand(initCmd)
}
