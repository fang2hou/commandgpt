package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"commandgpt/view"
)

func main() {
	v := view.New()
	p := tea.NewProgram(v, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(v.Result)
}
