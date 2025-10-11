package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"

	"commit-craft/config"
	"commit-craft/provider"

	"github.com/charmbracelet/huh"
)

func executeGitCommit(msg string) {
	log.Printf("Running command: git commit -m %s \n", msg)
	_, err := exec.Command("git", "commit", "-m", msg).Output()
	if err != nil {
		log.Fatalf("failed to execute git commit: %v", err)
	}
}

const (
	StateCommit = iota
	StateQuit
	StateModify
)

func showOutputScreen(msg string) {
	log.Printf("Generated commit msg: '%s' \n", msg)
	state := StateQuit
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose your action:").
				Options(
					huh.NewOption("Commit", StateCommit),
					huh.NewOption("Modify", StateModify),
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

func main() {
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
}
