package main

import (
	"fmt"
	"os"

	"github.com/jhoffmann/bookmark-manager/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "bookmark-manager",
		Short: "A beautiful TUI bookmark manager for folders",
		Long: `Bookmark Manager - A command-line tool for managing folder bookmarks with a beautiful terminal interface.

Use this tool to:
- Add folder bookmarks with custom categories
- Browse bookmarks with an interactive TUI
- Filter and search through your bookmarks
- Export bookmarks to JSON

Run without arguments to see this help, or use one of the subcommands.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Default to help when no subcommands
			cmd.Help()
		},
	}

	// Add subcommands
	rootCmd.AddCommand(cmd.GetAddCmd())
	rootCmd.AddCommand(cmd.GetListCmd())
	rootCmd.AddCommand(cmd.GetExportCmd())

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
