// Package styles provides consistent styling for the TUI using lipgloss.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for the bookmark manager TUI
var (
	Primary   = lipgloss.Color("63")  // Purple
	Secondary = lipgloss.Color("39")  // Light blue
	Success   = lipgloss.Color("42")  // Green
	Warning   = lipgloss.Color("214") // Orange
	Error     = lipgloss.Color("196") // Red
	Muted     = lipgloss.Color("245") // Gray
	Subtle    = lipgloss.Color("241") // Dark gray
	Border    = lipgloss.Color("238") // Border gray
)

// Base styles
var (
	// BaseStyle is the foundation for all components
	BaseStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// BorderStyle for containers
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	// FocusedBorderStyle for focused containers
	FocusedBorderStyle = BorderStyle.Copy().
				BorderForeground(Primary)
)

// Tab styles
var (
	// ActiveTab style for the currently selected tab
	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(Primary).
			Padding(0, 2).
			MarginRight(1)

	// InactiveTab style for non-selected tabs
	InactiveTab = lipgloss.NewStyle().
			Foreground(Muted).
			Background(lipgloss.Color("236")).
			Padding(0, 2).
			MarginRight(1)

	// TabContainer style for the tab bar
	TabContainer = lipgloss.NewStyle().
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottomForeground(Border).
			MarginBottom(1)
)

// List styles
var (
	// SelectedItem style for the currently highlighted list item
	SelectedItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(Primary).
			Bold(true).
			Padding(0, 1)

	// UnselectedItem style for non-highlighted list items
	UnselectedItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Padding(0, 1)

	// ItemPath style for the folder path in list items
	ItemPath = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Width(50)

	// ItemCategory style for the category badge in list items
	ItemCategory = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(Secondary).
			Padding(0, 1).
			MarginLeft(2).
			Bold(true)

	// ListContainer style for the list wrapper
	ListContainer = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Height(20).
			Padding(1)
)

// Input styles
var (
	// FilterInput style for the search filter
	FilterInput = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Muted).
			Padding(0, 1).
			MarginBottom(1)

	// FocusedFilterInput style when filter is focused
	FocusedFilterInput = FilterInput.Copy().
				BorderForeground(Primary)
)

// Status styles
var (
	// StatusBar style for the bottom status line
	StatusBar = lipgloss.NewStyle().
			Foreground(Muted).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTopForeground(Border).
			MarginTop(1).
			Padding(1, 2)

	// HelpText style for help information
	HelpText = lipgloss.NewStyle().
			Foreground(Subtle)

	// CountText style for showing bookmark counts
	CountText = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)
)

// Dialog styles
var (
	// DialogBox style for confirmation dialogs
	DialogBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(2, 4).
			Background(lipgloss.Color("235")).
			Width(50).
			Align(lipgloss.Center)

	// DialogTitle style for dialog titles
	DialogTitle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(2)

	// DialogText style for dialog content
	DialogText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Center).
			MarginBottom(1)

	// ConfirmButton style for confirmation buttons
	ConfirmButton = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(Error).
			Padding(0, 3).
			Bold(true).
			Border(lipgloss.RoundedBorder())

	// CancelButton style for cancel buttons
	CancelButton = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(Muted).
			Padding(0, 3).
			Bold(true).
			Border(lipgloss.RoundedBorder())
)

// Success/Error message styles
var (
	// SuccessMessage style for success notifications
	SuccessMessage = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true).
			Padding(0, 1)

	// ErrorMessage style for error notifications
	ErrorMessage = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true).
			Padding(0, 1)

	// WarningMessage style for warning notifications
	WarningMessage = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true).
			Padding(0, 1)
)

// Utility functions for dynamic styling

// GetCategoryStyle returns a style for category badges with different colors
func GetCategoryStyle(category string) lipgloss.Style {
	switch category {
	case "work":
		return ItemCategory.Copy().Background(lipgloss.Color("208")) // Orange
	case "personal":
		return ItemCategory.Copy().Background(lipgloss.Color("35")) // Green
	case "default":
		return ItemCategory.Copy().Background(Muted)
	default:
		return ItemCategory.Copy().Background(Secondary)
	}
}

// Center centers text within a given width
func Center(text string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(text)
}

// Truncate truncates text to fit within a given width with ellipsis
func Truncate(text string, width int) string {
	if lipgloss.Width(text) <= width {
		return text
	}
	if width < 3 {
		return text[:width]
	}
	return text[:width-3] + "..."
}
