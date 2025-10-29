package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"time"

	"commit-craft/config"
	"commit-craft/provider"
	"commit-craft/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"golang.design/x/clipboard"
)

func executeGitCommit(msg string) {
	log.Printf("Running command: git commit -m \"%s\"\n", msg)
	_, err := exec.Command("git", "commit", "-m", msg).Output()
	if err != nil {
		log.Fatalf("failed to execute git commit: %v", err)
	}
}

const (
	StateCommit = iota
	StateQuit
	StateModify
	StateCopyToClipboard
)

func copyToClipboard(msg string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtText, []byte(msg))
	// TODO: Technical Debt (If dont sleep, it can not save to clipboard)
	time.Sleep(time.Second)
	return nil
}

func showOutputScreen(msg string) {
	log.Printf("Generated commit msg:\n%s\n", msg)
	state := StateQuit
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose your action:").
				Options(
					huh.NewOption("Commit", StateCommit),
					huh.NewOption("Modify", StateModify),
					huh.NewOption("Copy to Clipboard", StateCopyToClipboard),
					huh.NewOption("Quit", StateQuit),
				).
				Value(&state)),
	).Run()
	if err != nil {
		log.Fatalln("Choose action error:", err)
	}
	switch state {
	case StateModify:
		err := huh.NewForm(huh.NewGroup(huh.NewText().Title("Edit your commit message:").Value(&msg))).Run()
		if err != nil {
			log.Fatal(err)
		}
		showOutputScreen(msg)
		return
	case StateCommit:
		executeGitCommit(msg)
		return
	case StateCopyToClipboard:
		err := copyToClipboard(msg)
		if err != nil {
			log.Fatal("Failed to copy to clipboard: ", err)
		}
		return
	case StateQuit:
		os.Exit(0)
	}
}

func executeGitDiff() (string, error) {
	diff, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		return "", errors.New("failed to execute git diff --cached" + err.Error())
	}
	if string(diff) == "" {
		return "", errors.New("no changes to commit, working tree clean")
	}
	return string(diff), nil
}

func temp() {
	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "")
		if err != nil {
			log.Fatalln("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close() // nolint:errcheck
	}
	model := tui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	temp()
	os.Exit(0)
	ctx := context.Background()
	config, err := config.BuildConfig()
	if err != nil {
		log.Fatalln(`failed to build commit-craft config`, err)
	}
	diff, err := executeGitDiff()
	if err != nil {
		log.Fatalln(`failed to execute "git diff --cached"`, err)
	}
	commitCraftClient, err := provider.GetClient(config)
	if err != nil {
		log.Fatalln(`failed to init commit craft client `, err)
	}
	result, err := commitCraftClient.GenerateCommit(ctx, diff)
	if err != nil {
		log.Fatalln("Failed to generate content " + err.Error())
	}
	msg := result
	showOutputScreen(msg)
	os.Exit(0)
}
