package tui

import (
	"log"
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
	StateLoadingDiff
	StateGenerating
	StateReview
	StateEdit
	StateAllDone
)

type ReviewAction int

const (
	ActionCommit ReviewAction = iota
	ActionModify
	ActionCopy
	ActionQuit
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
	Enter key.Binding
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
		key.WithKeys("q", "ctrl+c", "escalpe"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

type Model struct {
	currState    uint
	form         *huh.Form
	config       config.Config
	keys         keyMap
	help         help.Model
	diff         string
	commitMsg    string
	err          error
	spinner      bool // simple flag to show loading state if needed, or use separate bubble
	reviewAction ReviewAction
}

func getReviewForm(action *ReviewAction) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[ReviewAction]().
				Title("Choose your action:").
				Options(
					huh.NewOption("Commit", ActionCommit),
					huh.NewOption("Modify", ActionModify),
					huh.NewOption("Copy to Clipboard", ActionCopy),
					huh.NewOption("Quit", ActionQuit),
				).
				Value(action),
		),
	)
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
	log.Println("Initializing model")
	m.help.ShowAll = true
	return func() tea.Msg {
		log.Println("Checking API key environment variable")
		apiKey, isEnvAPIKeyFound := os.LookupEnv("COMMIT_CRAFT_GEMINI_KEY1")

		if !isEnvAPIKeyFound {
			log.Println("API key not found in env, triggering form")
			m.form = getSetAPIKeyForm(&m.config.APIKey)
			m.currState = InitState
			return m.form.Init()
		}

		log.Println("API key found in env")
		m.config.APIKey = apiKey
		// Instead of going to MainState immediately, check Diff or Config
		// For now we can transition to LoadingDiff via Update loop or return a specific Msg
		// But to keep it simple, let's trigger a CheckDiffMsg
		return CheckDiffMsg{}
	}
}

type CheckDiffMsg struct{}

func (m *Model) Init() tea.Cmd {
	return initModel(m)
}
