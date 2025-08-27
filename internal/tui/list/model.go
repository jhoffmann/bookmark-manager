// Package list provides the main TUI list interface for browsing bookmarks.
package list

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhoffmann/bookmark-manager/internal/bookmark"
	"github.com/jhoffmann/bookmark-manager/internal/tui/confirm"
	"github.com/jhoffmann/bookmark-manager/internal/tui/styles"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// Model represents the main TUI state for the bookmark list
type Model struct {
	list           list.Model
	categories     []string
	activeCategory string
	filter         textinput.Model
	filterFocused  bool
	bookmarks      []*bookmark.Bookmark
	allBookmarks   []*bookmark.Bookmark
	service        *bookmark.Service
	keys           keyMap
	confirmDialog  confirm.Model
	showingDialog  bool
	windowSize     tea.WindowSizeMsg
	err            error
}

// bookmarkItem implements list.Item for use with bubbles/list
type bookmarkItem struct {
	bookmark *bookmark.Bookmark
}

func (i bookmarkItem) FilterValue() string {
	return i.bookmark.Folder + " " + string(i.bookmark.Category)
}

func (i bookmarkItem) Title() string {
	return i.bookmark.Folder
}

func (i bookmarkItem) Description() string {
	return string(i.bookmark.Category)
}

// keyMap defines key bindings for the list interface
type keyMap struct {
	NextTab     key.Binding
	PrevTab     key.Binding
	Delete      key.Binding
	Filter      key.Binding
	Open        key.Binding
	Quit        key.Binding
	ClearFilter key.Binding
	Enter       key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() keyMap {
	return keyMap{
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next category"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev category"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x", "d"),
			key.WithHelp("x/d", "delete bookmark"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter bookmarks"),
		),
		Open: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open folder"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "clear filter"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open folder"),
		),
	}
}

// New creates a new list model
func New(service *bookmark.Service, initialCategory string) Model {
	// Initialize text input for filtering
	filterInput := textinput.New()
	filterInput.Placeholder = "Type to filter bookmarks..."
	filterInput.CharLimit = 156

	// Initialize list
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = "All" // Start with "All" category
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false) // We handle filtering ourselves

	return Model{
		list:           l,
		categories:     []string{"All"},
		activeCategory: initialCategory,
		filter:         filterInput,
		filterFocused:  false,
		keys:           DefaultKeyMap(),
		confirmDialog:  confirm.New(),
		service:        service,
	}
}

// Init initializes the model (required by tea.Model interface)
func (m Model) Init() tea.Cmd {
	return m.LoadBookmarks()
}

// LoadBookmarks loads bookmarks from the database
func (m *Model) LoadBookmarks() tea.Cmd {
	return func() tea.Msg {
		// Load all bookmarks
		allBookmarks, err := bookmark.List(m.service, 0, 0)
		if err != nil {
			return errMsg{err}
		}

		// Extract unique categories
		categorySet := make(map[string]bool)
		categorySet["All"] = true
		categorySet["default"] = true

		for _, b := range allBookmarks {
			categorySet[string(b.Category)] = true
		}

		categories := make([]string, 0, len(categorySet))
		categories = append(categories, "All")
		for cat := range categorySet {
			if cat != "All" {
				categories = append(categories, cat)
			}
		}
		sort.Strings(categories[1:]) // Sort all except "All"

		return bookmarksLoadedMsg{
			bookmarks:  allBookmarks,
			categories: categories,
		}
	}
}

// Update handles input events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle confirmation dialog first
	if m.showingDialog {
		newConfirm, confirmCmd := m.confirmDialog.Update(msg)
		m.confirmDialog = newConfirm

		// Check if user made a choice
		if m.confirmDialog.HasResult() {
			m.showingDialog = false
			result := m.confirmDialog.GetResult()
			if result.Confirmed && result.Bookmark != nil {
				return m, m.deleteBookmark(result.Bookmark)
			}
		}

		return m, confirmCmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg // Store the current window size
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4) // Reserve space for filter
		m.filter.Width = msg.Width - h - 10

	case tea.KeyMsg:
		// Handle filter input when focused
		if m.filterFocused {
			switch msg.String() {
			case "esc":
				m.filterFocused = false
				m.filter.Blur()
			case "enter":
				m.filterFocused = false
				m.filter.Blur()
				return m, m.applyFilter()
			default:
				m.filter, cmd = m.filter.Update(msg)
				cmds = append(cmds, cmd)
				// Auto-apply filter as user types
				cmds = append(cmds, m.applyFilter())
			}
			return m, tea.Batch(cmds...)
		}

		// Handle main interface key bindings
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Filter):
			m.filterFocused = true
			m.filter.Focus()

		case key.Matches(msg, m.keys.ClearFilter):
			m.filter.SetValue("")
			return m, m.applyFilter()

		case key.Matches(msg, m.keys.NextTab):
			return m, m.nextCategory()

		case key.Matches(msg, m.keys.PrevTab):
			return m, m.prevCategory()

		case key.Matches(msg, m.keys.Delete):
			if selectedItem, ok := m.list.SelectedItem().(bookmarkItem); ok {
				m.confirmDialog.Show(
					selectedItem.bookmark,
					"", // No longer needed with new simple dialog
				)
				m.showingDialog = true
				// Immediately send the current window size to the confirm dialog
				if m.windowSize.Width > 0 && m.windowSize.Height > 0 {
					m.confirmDialog, _ = m.confirmDialog.Update(m.windowSize)
				}
			}

		case key.Matches(msg, m.keys.Open, m.keys.Enter):
			if selectedItem, ok := m.list.SelectedItem().(bookmarkItem); ok {
				return m, m.openFolder(selectedItem.bookmark.Folder)
			}
		}

	case bookmarksLoadedMsg:
		m.allBookmarks = msg.bookmarks
		m.categories = msg.categories

		// Set initial category if not already set
		if m.activeCategory == "" {
			m.activeCategory = "All"
		}

		// Update list title to show current category
		m.list.Title = m.activeCategory

		return m, m.filterByCategory()

	case bookmarksFilteredMsg:
		m.bookmarks = msg.bookmarks

		// Convert to list items
		items := make([]list.Item, len(m.bookmarks))
		for i, b := range m.bookmarks {
			items[i] = bookmarkItem{bookmark: b}
		}

		m.list.SetItems(items)

	case bookmarkDeletedMsg:
		return m, m.LoadBookmarks() // Reload bookmarks

	case errMsg:
		m.err = msg.err
	}

	// Update list
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the interface
func (m Model) View() string {
	if m.showingDialog {
		return m.confirmDialog.View()
	}

	// Render filter if focused or has value
	var filterView string
	if m.filterFocused {
		filterStyle := styles.FocusedFilterInput
		filterView = filterStyle.Render("Filter: "+m.filter.View()) + "\n"
	} else if m.filter.Value() != "" {
		filterStyle := styles.FilterInput
		filterView = filterStyle.Render("Filter: "+m.filter.View()) + "\n"
	}

	return docStyle.Render(filterView + m.list.View())
}

// Helper functions

func (m *Model) nextCategory() tea.Cmd {
	for i, cat := range m.categories {
		if cat == m.activeCategory {
			nextIndex := (i + 1) % len(m.categories)
			m.activeCategory = m.categories[nextIndex]
			m.list.Title = m.activeCategory // Update list title
			break
		}
	}
	return m.filterByCategory()
}

func (m *Model) prevCategory() tea.Cmd {
	for i, cat := range m.categories {
		if cat == m.activeCategory {
			prevIndex := (i - 1 + len(m.categories)) % len(m.categories)
			m.activeCategory = m.categories[prevIndex]
			m.list.Title = m.activeCategory // Update list title
			break
		}
	}
	return m.filterByCategory()
}

func (m *Model) filterByCategory() tea.Cmd {
	return func() tea.Msg {
		var filtered []*bookmark.Bookmark

		if m.activeCategory == "All" {
			filtered = m.allBookmarks
		} else {
			for _, b := range m.allBookmarks {
				if string(b.Category) == m.activeCategory {
					filtered = append(filtered, b)
				}
			}
		}

		return bookmarksFilteredMsg{bookmarks: filtered}
	}
}

func (m *Model) applyFilter() tea.Cmd {
	return func() tea.Msg {
		filterText := strings.ToLower(strings.TrimSpace(m.filter.Value()))

		var filtered []*bookmark.Bookmark
		source := m.allBookmarks

		// First filter by category
		if m.activeCategory != "All" {
			source = []*bookmark.Bookmark{}
			for _, b := range m.allBookmarks {
				if string(b.Category) == m.activeCategory {
					source = append(source, b)
				}
			}
		}

		// Then apply text filter
		if filterText == "" {
			filtered = source
		} else {
			for _, b := range source {
				if strings.Contains(strings.ToLower(b.Folder), filterText) ||
					strings.Contains(strings.ToLower(string(b.Category)), filterText) {
					filtered = append(filtered, b)
				}
			}
		}

		return bookmarksFilteredMsg{bookmarks: filtered}
	}
}

func (m *Model) deleteBookmark(b *bookmark.Bookmark) tea.Cmd {
	return func() tea.Msg {
		if err := b.Delete(m.service); err != nil {
			return errMsg{err}
		}
		return bookmarkDeletedMsg{}
	}
}

func (m *Model) openFolder(path string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", path)
		case "windows":
			cmd = exec.Command("explorer", path)
		default: // linux and others
			cmd = exec.Command("xdg-open", path)
		}

		if err := cmd.Run(); err != nil {
			return errMsg{fmt.Errorf("failed to open folder: %w", err)}
		}

		return nil
	}
}

// Messages
type bookmarksLoadedMsg struct {
	bookmarks  []*bookmark.Bookmark
	categories []string
}

type bookmarksFilteredMsg struct {
	bookmarks []*bookmark.Bookmark
}

type bookmarkDeletedMsg struct{}

type errMsg struct {
	err error
}
