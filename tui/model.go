package tui

import (
	"log"
	"os"

	"commit-craft/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type Model struct {
	form   *huh.Form
	config config.Config
}

func GetSetAPIKeyForm(apiKey *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Can not find Gemini Api key in your env, what is your Gemini Api Key?").Value(apiKey).Key("apiKey"),
		),
	)
}

func NewModel() *Model {
	return &Model{}
}

func initModel(m *Model) tea.Cmd {
	log.Println("aklhjfaskl")
	return func() tea.Msg {
		apiKey, isEnvAPIKeyFound := os.LookupEnv("COMMIT_CRAFT_GEMINI_KEY")

		if !isEnvAPIKeyFound {
			m.form = GetSetAPIKeyForm(&m.config.APIKey)
			return tea.Batch(m.form.Init(), tea.EnterAltScreen)
		} else {
			m.config.APIKey = apiKey
		}

		return nil
	}
}

func (m *Model) Init() tea.Cmd {
	return initModel(m)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	if m.form != nil {

		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		return m, cmd
	}
	return m, nil
}
