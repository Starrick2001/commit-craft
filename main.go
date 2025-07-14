package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("COMMIT_CRAFT_GEMINI_KEY")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatal(err)
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text("Explain how AI works in a few words"),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: Int32(0), // Disables thinking
			},
		},
	)

	fmt.Println(result.Text())
}
