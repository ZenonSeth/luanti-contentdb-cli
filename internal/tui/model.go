// Package tui is the interactive terminal UI, built on bubbletea.
// Each screen has its own file (menu.go, mods.go, ...); this file
// holds the shared model and dispatches Update/View to the current
// screen.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ZenonSeth/luanti-mod-manager-cli/internal/content"
)

type screen int

const (
	screenMenu screen = iota
	screenMods
	screenModDetail
)

type model struct {
	screen  screen
	cursor  int
	choices []string
	modsDir string

	mods       []content.Mod
	modsErr    error
	modsCursor int

	selectedMod content.Mod
}

// New returns the initial TUI model, ready to pass to tea.NewProgram.
func New(modsDir string) model {
	return model{
		screen:  screenMenu,
		cursor:  0,
		choices: []string{"List installed mods", "Quit"},
		modsDir: modsDir,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if keyMsg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.screen {
	case screenMenu:
		return m.updateMenu(keyMsg)
	case screenMods:
		return m.updateMods(keyMsg)
	case screenModDetail:
		return m.updateModDetail(keyMsg)
	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case screenMods:
		return m.viewMods()
	case screenModDetail:
		return m.viewModDetail()
	default:
		return m.viewMenu()
	}
}
