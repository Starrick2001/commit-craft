package tui

import (
	"os"

	"commit-craft/config"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

const (
	InitState uint = iota
	MainState
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

type Model struct {
	currState uint
	form      *huh.Form
	config    config.Config
	keys      keyMap
	help      help.Model
}

func getSetAPIKeyForm(apiKey *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Can not find Gemini Api key in your env, what is your Gemini Api Key?").Value(apiKey),
		),
	)
}

func NewModel() *Model {
	return &Model{currState: InitState, keys: keys, help: help.New()}
}

func initModel(m *Model) tea.Cmd {
	m.help.ShowAll = true
	return func() tea.Msg {
		apiKey, isEnvAPIKeyFound := os.LookupEnv("COMMIT_CRAFT_GEMINI_KEY1")

		if !isEnvAPIKeyFound {
			m.form = getSetAPIKeyForm(&m.config.APIKey)
			m.currState = InitState
			return m.form.Init()
		} else {
			m.currState = MainState
			m.config.APIKey = apiKey
		}

		return nil
	}
}

func (m *Model) Init() tea.Cmd {
	return initModel(m)
}
