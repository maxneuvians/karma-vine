package main

import (
	"fmt"
	"os"

	"karma_vine/internal/game"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(game.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
