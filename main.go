package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"google.golang.org/genai"
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
		os.Exit(3)
	}
}

func main() {
	ctx := context.Background()
	apiKey, isEnvApiKeyFound := os.LookupEnv("COMMIT_CRAFT_GEMINI_KEY")
	if !isEnvApiKeyFound {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Can not find Gemini Api key in your env, what is your Gemini Api Key?").Value(&apiKey),
			),
		)
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}
	}
	diff, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		log.Fatalln("failed to execute git diff --cached", err)
	}
	if string(diff) == "" {
		log.Fatalln("no changes to commit, working tree clean")
	}
	config := Config{apiKey: apiKey}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a model:").
				Options(
					huh.NewOption("gemini-2.0-flash", "gemini-2.0-flash"),
					huh.NewOption("gemini-2.5-flash", "gemini-2.5-flash")).
				Value(&config.model)),
	)
	if err := form.Run(); err != nil {
		log.Fatal(err)
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: config.apiKey})
	if err != nil {
		log.Fatalln("Failed to create Gemini client:%w", err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text("You are an expert at creating a git commit message for a set of changes. Return the generated title commit message. Here is a diff of changes we need a commit message for: "+string(diff)),
		// genai.Text("You are an expert at creating a git commit message for a set of changes. Return a git commit command line with generated commit message. Here is a diff of changes we need a commit message for: "+string(diff)),
		&genai.GenerateContentConfig{
			// ThinkingConfig: &genai.ThinkingConfig{
			// ThinkingBudget: Int32(0), // Disables thinking
			// },
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	msg := result.Text()
	showOutputScreen(msg)
}
