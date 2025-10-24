package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.form != nil {
		if m.form.State == huh.StateCompleted {
			apiKey := m.form.GetString("apiKey")
			return fmt.Sprintf("You selected: %s\n\n", apiKey) + lipgloss.NewStyle().Background(lipgloss.Color("99")).Align(lipgloss.Center).Render("ctrl+c - quit")

		}
		return m.form.View()
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render("ctrl+c - quit")
}
