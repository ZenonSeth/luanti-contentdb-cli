package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ZenonSeth/luanti-mod-manager-cli/internal/tui"
)

func main() {
	dir := flag.String("dir", ".", "path to the mods folder")
	flag.Parse()

	p := tea.NewProgram(tui.New(*dir))
	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
