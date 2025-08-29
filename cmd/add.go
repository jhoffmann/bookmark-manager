package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jhoffmann/bookmark-manager/internal/app"
	"github.com/jhoffmann/bookmark-manager/internal/models"
	"github.com/jhoffmann/bookmark-manager/internal/tui/styles"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [category]",
	Short: "Add the current directory as a bookmark",
	Long: `Add the current directory as a bookmark with an optional category.

Examples:
  bookmark-manager add
  bookmark-manager add work
  bookmark-manager add personal
  bookmark-manager add "my-project"`,
	Args: cobra.MaximumNArgs(1),
	Run:  runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	// Initialize app (loads config, database, and service)
	appInstance := app.InitializeOrExit()
	defer appInstance.Close()

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s Failed to get current directory: %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}

	// Get absolute path to ensure consistency
	absPath, err := filepath.Abs(currentDir)
	if err != nil {
		fmt.Printf("%s Failed to get absolute path: %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}

	// Determine category
	var category models.CategoryType
	if len(args) > 0 && args[0] != "" {
		category = models.CategoryType(args[0])
	}
	// If no category provided, it will remain empty

	// Check if bookmark already exists
	existingBookmarks, err := appInstance.Service.SearchByFolder(absPath)
	if err != nil {
		fmt.Printf("%s Failed to check for existing bookmarks: %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}

	// Check for exact match
	for _, existing := range existingBookmarks {
		if existing.Folder == absPath {
			fmt.Printf("%s Bookmark already exists: %s [%s]\n",
				styles.WarningMessage.Render("!"),
				existing.Folder,
				existing.Category)
			return
		}
	}

	// Create new bookmark
	newBookmark := &models.Bookmark{
		Folder:   absPath,
		Category: category,
	}

	// Save bookmark
	if err := appInstance.Service.Save(newBookmark); err != nil {
		fmt.Printf("%s Failed to save bookmark: %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}

	// Success message
	fmt.Printf("%s Added bookmark: %s [%s]\n",
		styles.SuccessMessage.Render("✓"),
		absPath,
		category)
}

// GetAddCmd returns the add command
func GetAddCmd() *cobra.Command {
	return addCmd
}

func init() {
	// Add flags if needed in the future
	// addCmd.Flags().BoolVarP(&force, "force", "f", false, "Force add even if bookmark exists")
}
