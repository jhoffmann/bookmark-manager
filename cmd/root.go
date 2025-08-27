package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bookmark-manager",
	Short: "A command line bookmark manager",
	Long:  "A command line utility for managing bookmarks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to bookmark-manager!")
		fmt.Println("Use --help to see available commands")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
