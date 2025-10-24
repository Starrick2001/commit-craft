package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func renderInitScreen(m *Model) string {
	if m.form != nil {
		return m.form.View()
	}
	return "Initializing..."
}

func renderMainScreen(m *Model) string {
	s := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(fmt.Sprintf("Your apiKey: %s", (m.config.APIKey)))
	return s
}

func renderLoading(m *Model, title string) string {
	return fmt.Sprintf("\n %s\n", title) // TODO: Use m.spinner here
}

func renderReview(m *Model) string {
	s := fmt.Sprintf("Generated commit msg:\n%s\n\n", m.commitMsg)
	if m.form != nil {
		s += m.form.View()
	}
	return s
}

func (m *Model) View() string {
	s := ""
	switch m.currState {
	case InitState:
		return renderInitScreen(m)
	case MainState:
		s += renderMainScreen(m)
	case StateLoadingDiff:
		s += renderLoading(m, "Checking git diff...")
	case StateGenerating:
		s += renderLoading(m, "Generating commit message...")
	case StateReview:
		s += renderReview(m)
	case StateEdit:
		if m.form != nil {
			s += m.form.View()
		}
	case StateAllDone:
		s += "Done!\n"
	}

	if m.err != nil {
		s += fmt.Sprintf("\nError: %v\n", m.err)
		s += "\nPress q to quit.\n"
	}

	s += "\n"
	if m.currState != InitState && m.currState != StateReview && m.currState != StateEdit {
		// TODO: Bring back help view
		// Only show help keys if not in a form (forms handle their own keys usually)
		// But basic help is useful.
		// m.help.View(m.keys) might conflict with form keys if not careful.
		// For now, minimal help.
	}
	return s
}
