package main

import (
	"fmt"
	"os"

	"changelog-generator/internal/lib"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration from .changelogrc.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		//Load the config file
		config, err := lib.LoadConfig(".changelogrc.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println()
			fmt.Println("Tip: Run 'changelog init' to create a config file")
			os.Exit(1)
		}
		//Display the config
		lib.PrintConfig(config)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
