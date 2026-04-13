package main

import (
	"fmt"
	"os"

	"karma_vine/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(game.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
