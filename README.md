# 📁 Bookmark Manager

A beautiful command-line bookmark manager for folders with an interactive Terminal User Interface (TUI) built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

![demo](https://github.com/jhoffmann/bookmark-manager/blob/main/demo/demo.gif?raw=true)

![Go Version](https://img.shields.io/badge/Go-1.24+-blue)
![License](https://img.shields.io/github/license/jhoffmann/bookmark-manager)
![Commits](https://img.shields.io/github/commit-activity/t/jhoffmann/bookmark-manager)

## ✨ Features

- 🎯 **Beautiful TUI**: Interactive terminal interface with syntax highlighting and smooth navigation
- 📂 **Folder Bookmarks**: Bookmark any folder on your system, not just URLs
- 🏷️ **Custom Categories**: Organize bookmarks with user-defined categories
- 🔍 **Smart Filtering**: Real-time search and filtering capabilities
- ⌨️ **Keyboard Navigation**: Full keyboard shortcuts for efficient workflow
- 📤 **JSON Export**: Export bookmarks in JSON format for backup or integration
- 🖥️ **Cross-Platform**: Works on Linux, macOS, and Windows
- 🗃️ **SQLite Storage**: Reliable local database storage
- 🎨 **Modern Design**: Styled with lipgloss for a polished look

## 🚀 Installation

### Build from Source

```bash
git clone https://github.com/jhoffmann/bookmark-manager.git
cd bookmark-manager
CGO_ENABLED=1 go build -o bookmark-manager .
```

## 📖 Usage

### Basic Commands

```bash
# Show help (default when no arguments provided)
./bookmark-manager

# Add current directory as bookmark
./bookmark-manager add [category]

# Launch interactive TUI browser
./bookmark-manager list [category]

# Export bookmarks to JSON
./bookmark-manager export [category] [filter]
```

### Examples

```bash
# Add the current folder
./bookmark-manager add

# Add to a custom category
./bookmark-manager add "personal"

# Launch TUI showing all bookmarks
./bookmark-manager list

# Launch TUI filtered to "personal" category
./bookmark-manager list personal

# Export all bookmarks to JSON
./bookmark-manager export > my-bookmarks.json

# Export only work bookmarks
./bookmark-manager export personal > personal-bookmarks.json

# Export filtered bookmarks
./bookmark-manager export "" projects > project-bookmarks.json
```

### Shell Integration

Here is a handy shell function to switch the current directory rather than opening a new terminal/explorer:

```bash
function bm() {
    local tmp="$(mktemp -t "bookmark-manager-cwd.XXXXXX")"
    bookmark-manager list "$@" --cwd-file="$tmp"
    if cwd="$(cat -- "$tmp")" && [ -n "$cwd" ] && [ "$cwd" != "$PWD" ]; then
        cd -- "$cwd"
    fi
    rm -f -- "$tmp"
}

# See bookmark-manager completion for your shell
eval "$(bookmark-manager completion bash)"

```

## ⚙️ Configuration

### Cross-Platform Database Locations

- **Linux**: `~/.config/bookmark-manager/bookmarks.db`
- **macOS**: `~/Library/Application Support/bookmark-manager/bookmarks.db`
- **Windows**: `%APPDATA%/bookmark-manager/bookmarks.db`

The application automatically creates the directory if it doesn't exist.

## 📊 JSON Export Format

```json
[
  {
    "id": 1,
    "folder": "/home/user/projects/awesome-project",
    "category": "work",
    "date_created": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "folder": "/home/user/documents/personal",
    "category": "personal",
    "date_created": "2024-01-15T11:15:00Z"
  }
]
```

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the amazing Bubble Tea ecosystem
- [GORM](https://gorm.io/) for excellent ORM capabilities
- [Cobra](https://cobra.dev/) for powerful CLI framework
- The Go community for excellent tooling and libraries

---

**Built with ❤️ using Go and Bubble Tea** 🫧🍃
