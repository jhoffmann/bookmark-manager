// Package confirm provides a simple list-based confirmation for dangerous operations.
package confirm

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhoffmann/bookmark-manager/internal/bookmark"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// choiceItem implements list.Item for yes/no choices
type choiceItem struct {
	title       string
	description string
	value       bool
}

func (i choiceItem) FilterValue() string { return i.title }
func (i choiceItem) Title() string       { return i.title }
func (i choiceItem) Description() string { return i.description }

// Model represents the confirmation state
type Model struct {
	list     list.Model
	bookmark *bookmark.Bookmark
	visible  bool
	result   bool
	chosen   bool
}

// New creates a new confirmation model
func New() Model {
	items := []list.Item{
		choiceItem{title: "No", description: "Cancel - don't delete", value: false},
		choiceItem{title: "Yes", description: "Delete this bookmark", value: true},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true

	l := list.New(items, delegate, 0, 0) // Use same sizing as main list
	l.Title = "Delete Bookmark?"
	l.SetShowStatusBar(true) // Same as main list
	l.SetFilteringEnabled(false)

	return Model{
		list:    l,
		visible: false,
		chosen:  false,
	}
}

// Show displays the confirmation with the given bookmark
func (m *Model) Show(bookmark *bookmark.Bookmark, message string) {
	m.bookmark = bookmark
	m.visible = true
	m.chosen = false
	m.result = false
	m.list.Select(0) // Default to "No"

	// Update title to include bookmark info
	if bookmark != nil {
		m.list.Title = "Delete: " + bookmark.Folder + "?"
	}
}

// Hide hides the confirmation
func (m *Model) Hide() {
	m.visible = false
	m.bookmark = nil
	m.chosen = false
}

// IsVisible returns whether the confirmation is currently visible
func (m Model) IsVisible() bool {
	return m.visible
}

// Update handles input events for the confirmation
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Get selected choice
			if selectedItem, ok := m.list.SelectedItem().(choiceItem); ok {
				m.result = selectedItem.value
				m.chosen = true
				m.visible = false
				return m, nil
			}
		case "esc", "q", "n":
			m.result = false
			m.chosen = true
			m.visible = false
			return m, nil
		case "y":
			m.result = true
			m.chosen = true
			m.visible = false
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the confirmation
func (m Model) View() string {
	if !m.visible {
		return ""
	}
	return docStyle.Render(m.list.View())
}

// Result represents the result of the confirmation
type Result struct {
	Confirmed bool
	Bookmark  *bookmark.Bookmark
}

// GetResult returns the result based on current state
func (m Model) GetResult() Result {
	return Result{
		Confirmed: m.result,
		Bookmark:  m.bookmark,
	}
}

// HasResult returns whether the user has made a choice
func (m Model) HasResult() bool {
	return m.chosen
}
