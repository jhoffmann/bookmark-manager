package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jhoffmann/bookmark-manager/internal/app"
	"github.com/jhoffmann/bookmark-manager/internal/tui/list"
	"github.com/jhoffmann/bookmark-manager/internal/tui/styles"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [category]",
	Short: "Interactive TUI for browsing bookmarks",
	Long: `Launch an interactive TUI to browse, filter, and manage your bookmarks.

Features:
- Tab through categories (All and custom categories)
- Real-time filtering with '/' key
- Delete bookmarks with 'x' key (with confirmation)
- Open folders with 'o' or 'enter' key
- Full keyboard navigation

Examples:
  bookmark-manager list
  bookmark-manager list work
  bookmark-manager list personal`,
	Args: cobra.MaximumNArgs(1),
	Run:  runList,
}

func runList(cmd *cobra.Command, args []string) {
	// Initialize app (loads config, database, and service)
	appInstance := app.InitializeOrExit()
	defer appInstance.Close()

	// Parse arguments
	var initialCategory string

	if len(args) > 0 && args[0] != "" {
		initialCategory = args[0]
	} else {
		initialCategory = "All"
	}

	// Create TUI model
	model := list.New(appInstance.Service, initialCategory)

	// Create Bubble Tea program
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	finalModel, err := program.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s TUI error: %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}

	// Handle any final state from the model
	if listModel, ok := finalModel.(list.Model); ok {
		_ = listModel // We could handle final state here if needed
	}
}

// GetListCmd returns the list command
func GetListCmd() *cobra.Command {
	return listCmd
}
