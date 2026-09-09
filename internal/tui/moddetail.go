package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateModDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.screen = screenMods
	}

	return m, nil
}

func (m model) viewModDetail() string {
	mod := m.selectedMod

	s := fmt.Sprintf("%s\n\n", mod.Name)

	if mod.Title != "" {
		s += fmt.Sprintf("title:       %s\n", mod.Title)
	}
	if mod.Author != "" {
		s += fmt.Sprintf("author:      %s\n", mod.Author)
	}
	if mod.Description != "" {
		s += fmt.Sprintf("description: %s\n", mod.Description)
	}

	s += fmt.Sprintf("folder:      %s\n", mod.Dir)
	s += fmt.Sprintf("path:        %s\n", mod.Path)
	s += fmt.Sprintf("mod.conf ok: %v\n", mod.ConfOK)

	if len(mod.Depends) > 0 {
		s += fmt.Sprintf("depends:     %s\n", strings.Join(mod.Depends, ", "))
	}
	if len(mod.OptionalDepends) > 0 {
		s += fmt.Sprintf("optional:    %s\n", strings.Join(mod.OptionalDepends, ", "))
	}

	s += "\n(q to go back)\n"

	return s
}
