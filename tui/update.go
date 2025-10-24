package tui

import (
	"commit-craft/pkg/clipboard"
	"commit-craft/pkg/git"
	"commit-craft/provider"
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case CheckDiffMsg:
		log.Println("Checking git diff...")
		m.currState = StateLoadingDiff
		return m, func() tea.Msg {
			diff, err := git.Diff()
			if err != nil {
				log.Println("Error checking diff:", err)
				return ErrorMsg{err}
			}
			log.Println("Diff found, length:", len(diff))
			return DiffLoadedMsg{Diff: diff}
		}
	case DiffLoadedMsg:
		m.diff = msg.Diff
		m.currState = StateGenerating
		log.Println("Generating commit message...")
		return m, func() tea.Msg {
			ctx := context.Background()
			client, err := provider.GetClient(&m.config)
			if err != nil {
				log.Println("Error init client:", err)
				return ErrorMsg{err}
			}
			result, err := client.GenerateCommit(ctx, m.diff)
			if err != nil {
				log.Println("Error generating commit:", err)
				return ErrorMsg{err}
			}
			log.Println("Commit generated successfully")
			return CommitGeneratedMsg{Msg: result}
		}
	case CommitGeneratedMsg:
		m.commitMsg = msg.Msg
		m.currState = StateReview
		m.form = getReviewForm(&m.reviewAction)
		log.Println("Entering review state")
		return m, m.form.Init()
	case ErrorMsg:
		m.err = msg.Err
		log.Println("Error occurred:", msg.Err)
		return m, tea.Quit // Or show error state
	}

	switch m.currState {
	case InitState:
		if m.form != nil {
			form, cmd := m.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				m.form = f
			}
			cmds = append(cmds, cmd)

			if m.form.State == huh.StateCompleted {
				log.Println("Init form completed")
				m.currState = StateLoadingDiff
				cmds = append(cmds, func() tea.Msg { return CheckDiffMsg{} })
			}
		}
	case StateReview:
		if m.form != nil {
			form, cmd := m.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				m.form = f
			}
			cmds = append(cmds, cmd)

			if m.form.State == huh.StateCompleted {
				switch m.reviewAction {
				case ActionCommit:
					log.Println("Action: Commit")
					git.Commit(m.commitMsg)
					return m, tea.Quit
				case ActionModify:
					log.Println("Action: Modify")
					m.form = huh.NewForm(huh.NewGroup(huh.NewText().Title("Edit your commit message:").Value(&m.commitMsg)))
					m.currState = StateEdit
					return m, m.form.Init()
				case ActionCopy:
					log.Println("Action: Copy")
					clipboard.Copy(m.commitMsg)
					return m, tea.Quit
				case ActionQuit:
					log.Println("Action: Quit")
					return m, tea.Quit
				}
			}
		}
	case StateEdit:
		if m.form != nil {
			form, cmd := m.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				m.form = f
			}
			cmds = append(cmds, cmd)

			if m.form.State == huh.StateCompleted {
				log.Println("Edit completed")
				// Finished editing, go back to Review
				m.currState = StateReview
				m.form = getReviewForm(&m.reviewAction)
				cmds = append(cmds, m.form.Init())
			}
		}
	}
	return m, tea.Batch(cmds...)
}

type DiffLoadedMsg struct{ Diff string }
type CommitGeneratedMsg struct{ Msg string }
type ErrorMsg struct{ Err error }
