// Package styles provides consistent styling for the TUI using lipgloss.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for the bookmark manager TUI
var (
	Primary = lipgloss.Color("63")  // Purple
	Success = lipgloss.Color("42")  // Green
	Warning = lipgloss.Color("214") // Orange
	Error   = lipgloss.Color("196") // Red
	Muted   = lipgloss.Color("245") // Gray
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
