package main

import (
	"fmt"
	"os"

	"changelog-generator/internal/lib"

	"github.com/spf13/cobra"
)

// Flags for generate command
var (
	generateSince string
	generateTo    string
	commitCount   int
	outputFile    string
	useAI         bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a changelog",
	Long: `Analyze git commits and generate a changelog.

This command will:
  1. Load your configuration
  2. Scan your git repository
  3. Find commits between the specified range
  4. Generate a beautiful changelog`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		config, err := lib.LoadConfig(".changelogrc.yaml")
		if err != nil {
			fmt.Printf(" Error loading config: %v\n", err)
			fmt.Println()
			fmt.Println(" Tip: Run 'changelog init' to create a config file")
			os.Exit(1)
		}

		fmt.Println(" Generating changelog...")
		fmt.Println()
		fmt.Printf(" Project: %s\n", config.Project.Name)
		fmt.Printf("Repository: %s\n", config.Git.RepositoryPath)
		fmt.Println()

		// Open the git repository
		repo, err := lib.OpenRepository(config.Git.RepositoryPath)
		if err != nil {
			fmt.Printf(" Error opening repository: %v\n", err)
			fmt.Println()
			fmt.Println(" Make sure you're in a git repository!")
			os.Exit(1)
		}

		fmt.Printf(" Opened repository at: %s\n", config.Git.RepositoryPath)
		fmt.Println()

		// Get recent commits
		fmt.Printf(" Fetching last %d commits...\n", commitCount)
		commits, err := lib.GetRecentCommits(repo, commitCount)
		if err != nil {
			fmt.Printf(" Error getting commits: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf(" Found %d commits\n", len(commits))
		fmt.Println()

		//AI processing
		if useAI {
			aiClient, err := lib.NewAIClient()
			if err != nil {
				fmt.Printf(" AI not available: %v\n", err)
				fmt.Println("   Continuing without AI enhancement...")
				fmt.Println()
			} else {
				aiClient.ImproveAllCommits(commits)
			}
		}

		// Categorize and group commits
		groups := lib.GroupCommitsByCategory(commits)

		// Display grouped commits
		lib.PrintGroupedCommits(groups)

		// Generate markdown
		fmt.Println(" Generating markdown...")
		markdown := GenerateMarkdown(commits, config.Project.Name, config.Project.Version)

		// Determine output filename
		filename := outputFile
		if filename == "" {
			filename = config.Output.Filename
		}

		// Save to file
		err = SaveMarkdown(markdown, filename)
		if err != nil {
			fmt.Printf(" Error saving file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Changelog saved to: %s\n", filename)
		fmt.Println()
		fmt.Println(" Done!")
	},
}

func init() {
	// Add generate command to root command
	rootCmd.AddCommand(generateCmd)

	// Add flags to generate command
	generateCmd.Flags().StringVar(&generateSince, "since", "HEAD~10", "Starting point for changelog")
	generateCmd.Flags().StringVar(&generateTo, "to", "HEAD", "Ending point for changelog")
	generateCmd.Flags().IntVar(&commitCount, "count", 10, "Number of commits to show")
	generateCmd.Flags().StringVar(&outputFile, "output", "", "Output file (default from config)")
	generateCmd.Flags().BoolVar(&useAI, "ai", false, "Use AI to improve commit messages")

}
