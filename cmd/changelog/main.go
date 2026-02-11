package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "changelog",
	Short: "AI powered changelog generator",
	Long: `Changelog Generator - Create beautiful release notes automatically.
	This tool analyzes your git commits and generates professional changelogs
	using AI to understand what actually changed.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		// This runs when you just type "changelog" with no subcommands
		fmt.Println(" Welcome to Changelog Generator!")
		fmt.Println()
		fmt.Println("Available commands:")
		fmt.Println("  init      - Initialize configuration")
		fmt.Println("  generate  - Generate a changelog")
		fmt.Println()
		fmt.Println("Run 'changelog --help' for more information")
	},
}

func main() {
	//Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
	fmt.Println("Hello, world!")

}
