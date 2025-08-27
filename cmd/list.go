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
	// Get flag values
	cwdFile, _ := cmd.Flags().GetString("cwd-file")

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

	// Set cwd file mode if flag is provided
	if cwdFile != "" {
		model.SetCwdFile(cwdFile)
	}

	// Create Bubble Tea program
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	finalModel, err := program.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s TUI error: %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}

	// Handle any final state from the model - no additional handling needed
	// as the model writes to the file directly when a selection is made
	_ = finalModel
}

// GetListCmd returns the list command
func GetListCmd() *cobra.Command {
	listCmd.Flags().String("cwd-file", "", "Write the selection to the specified file and exit")
	return listCmd
}
