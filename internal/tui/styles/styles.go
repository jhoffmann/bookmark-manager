// Package styles provides consistent styling for the TUI using lipgloss.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for the bookmark manager TUI
var (
	Primary = lipgloss.Color("#e49fdb") // Magenta
	Success = lipgloss.Color("#7ac68f") // Green
	Warning = lipgloss.Color("#cdb36b") // Yellow
	Error   = lipgloss.Color("#fc9c93") // Red
	Muted   = lipgloss.Color("#404040") // Gray
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

// Message styles
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
