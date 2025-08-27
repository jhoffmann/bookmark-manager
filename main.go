package main

import (
	"fmt"
	"os"

	"github.com/jhoffmann/bookmark-manager/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// Get subcommands
	addCmd := cmd.GetAddCmd()
	listCmd := cmd.GetListCmd()
	exportCmd := cmd.GetExportCmd()

	rootCmd := &cobra.Command{
		Use:   "bookmark-manager",
		Short: "A beautiful TUI bookmark manager for folders",
		Run: func(rootCmd *cobra.Command, args []string) {
			// Default to list command when no subcommands
			listCmd.Run(rootCmd, args)
		},
	}

	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(exportCmd)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
