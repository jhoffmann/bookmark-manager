package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jhoffmann/bookmark-manager/internal/app"
	"github.com/jhoffmann/bookmark-manager/internal/bookmark"
	"github.com/jhoffmann/bookmark-manager/internal/tui/styles"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [category] [filter]",
	Short: "Export bookmarks to JSON",
	Long: `Export bookmarks to JSON format. Output is written to stdout for piping.

Examples:
  bookmark-manager export > all-bookmarks.json
  bookmark-manager export work > work-bookmarks.json
  bookmark-manager export personal home > personal-home-bookmarks.json
  bookmark-manager export "" projects > project-bookmarks.json`,
	Args: cobra.MaximumNArgs(2),
	Run:  runExport,
}

// ExportBookmark represents the JSON structure for exported bookmarks
type ExportBookmark struct {
	ID          uint   `json:"id"`
	Folder      string `json:"folder"`
	Category    string `json:"category"`
	DateCreated string `json:"date_created"`
}

func runExport(cmd *cobra.Command, args []string) {
	// Initialize app (loads config, database, and service)
	appInstance := app.InitializeOrExit()
	defer appInstance.Close()

	// Parse arguments
	var category bookmark.CategoryType
	var filter string

	if len(args) > 0 && args[0] != "" {
		category = bookmark.CategoryType(args[0])
	}

	if len(args) > 1 {
		filter = args[1]
	}

	// Fetch bookmarks based on criteria
	var bookmarks []*bookmark.Bookmark
	var err error

	if category != "" {
		// Search by category
		bookmarks, err = bookmark.SearchByCategory(appInstance.Service, category)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Failed to search bookmarks by category: %v\n",
				styles.ErrorMessage.Render("✗"), err)
			os.Exit(1)
		}
	} else {
		// Get all bookmarks
		bookmarks, err = bookmark.List(appInstance.Service, 0, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Failed to list bookmarks: %v\n",
				styles.ErrorMessage.Render("✗"), err)
			os.Exit(1)
		}
	}

	// Apply filter if specified
	if filter != "" {
		filteredBookmarks := make([]*bookmark.Bookmark, 0)
		filterLower := strings.ToLower(filter)

		for _, b := range bookmarks {
			if strings.Contains(strings.ToLower(b.Folder), filterLower) ||
				strings.Contains(strings.ToLower(string(b.Category)), filterLower) {
				filteredBookmarks = append(filteredBookmarks, b)
			}
		}
		bookmarks = filteredBookmarks
	}

	// Convert to export format
	exportBookmarks := make([]ExportBookmark, len(bookmarks))
	for i, b := range bookmarks {
		exportBookmarks[i] = ExportBookmark{
			ID:          b.ID,
			Folder:      b.Folder,
			Category:    string(b.Category),
			DateCreated: b.DateCreated.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Output JSON to stdout
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ") // Pretty print

	if err := encoder.Encode(exportBookmarks); err != nil {
		fmt.Fprintf(os.Stderr, "%s Failed to encode JSON: %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}
}

// GetExportCmd returns the export command
func GetExportCmd() *cobra.Command {
	return exportCmd
}
