package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func renderInitScreen(m *Model) string {
	if m.form != nil {
		return m.form.View()
	}
	return "m.form is nil"
}

func renderMainScreen(m *Model) string {
	string := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(fmt.Sprintf("Your apiKey: %s", (m.config.APIKey)))
	return string
}

func (m *Model) View() string {
	string := ""
	switch m.currState {
	case InitState:
		return renderInitScreen(m)
	case MainState:
		string += renderMainScreen(m)
	}
	string += "\n\n"
	string += m.help.View(m.keys)
	return string
}
