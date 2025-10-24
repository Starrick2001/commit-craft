package main

import (
	"log"
	"os"

	"commit-craft/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "")
		if err != nil {
			log.Fatalln("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	model := tui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
