package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"google.golang.org/genai"
)

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("COMMIT_CRAFT_GEMINI_KEY")
	diff, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		log.Fatalln(err)
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalln(err)
	}

	result, _ := client.Models.GenerateContent(
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

	fmt.Println("Generated commit msg: ", result.Text())
	exec.Command("git", "commit", "-m", result.Text()).Output()
}
