// Package edit provides a text input interface for editing bookmark categories.
package edit

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhoffmann/bookmark-manager/internal/bookmark"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

// Model represents the category editing state
type Model struct {
	textInput textinput.Model
	bookmark  *bookmark.Bookmark
	visible   bool
	result    string
	submitted bool
	cancelled bool
}

// New creates a new category edit model
func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter category name..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		textInput: ti,
		visible:   false,
		submitted: false,
		cancelled: false,
	}
}

// Show displays the edit dialog with the given bookmark
func (m *Model) Show(bookmark *bookmark.Bookmark) {
	m.bookmark = bookmark
	m.visible = true
	m.submitted = false
	m.cancelled = false
	m.result = ""

	// Pre-populate with current category
	currentCategory := string(bookmark.Category)
	m.textInput.SetValue(currentCategory)
	m.textInput.Focus()

	// Position cursor at end so user can easily edit
	m.textInput.CursorEnd()
}

// Hide hides the edit dialog
func (m *Model) Hide() {
	m.visible = false
	m.bookmark = nil
	m.submitted = false
	m.cancelled = false
	m.textInput.SetValue("")
}

// IsVisible returns whether the edit dialog is currently visible
func (m Model) IsVisible() bool {
	return m.visible
}

// Update handles input events for the edit dialog
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Submit the new category
			m.result = strings.TrimSpace(m.textInput.Value())
			m.submitted = true
			m.visible = false
			return m, nil
		case "esc":
			// Cancel editing
			m.cancelled = true
			m.visible = false
			return m, nil
		case "ctrl+c":
			// Also cancel on ctrl+c
			m.cancelled = true
			m.visible = false
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the edit dialog
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	title := "Edit Category"
	if m.bookmark != nil {
		title = "Edit Category: " + m.bookmark.Folder
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(title),
		"",
		m.textInput.View(),
		"",
		"Press Enter to save, Esc to cancel",
	)

	return docStyle.Render(content)
}

// Result represents the result of the category edit
type Result struct {
	NewCategory string
	Bookmark    *bookmark.Bookmark
	Submitted   bool
	Cancelled   bool
}

// GetResult returns the result based on current state
func (m Model) GetResult() Result {
	return Result{
		NewCategory: m.result,
		Bookmark:    m.bookmark,
		Submitted:   m.submitted,
		Cancelled:   m.cancelled,
	}
}

// HasResult returns whether the user has made a choice (submitted or cancelled)
func (m Model) HasResult() bool {
	return m.submitted || m.cancelled
}
