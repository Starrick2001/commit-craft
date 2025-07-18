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
	if _, err := exec.Command("git", "commit", "-m", msg).Output(); err != nil {
		log.Fatalf("failed to execute git commit: %v", err)
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
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
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

	log.Printf("Generated commit msg: %s \n", result.Text())
	confirm := false
	if err := huh.NewConfirm().Title("Do you want to exec git commit command?").Affirmative("Yes").Negative("No").Value(&confirm).Run(); err != nil {
		log.Fatalln("Confirmation error:", err)
	}
	if confirm {
		executeGitCommit(result.Text())
	}
}
