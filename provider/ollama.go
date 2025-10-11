package provider

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"commit-craft/config"
	"commit-craft/util"

	ollamaApi "github.com/ollama/ollama/api"
)

type Ollama struct {
	Config *config.Config
	client *ollamaApi.Client
}

func (o *Ollama) InitClient(ctx context.Context) error {
	u, err := url.Parse("http://localhost:11434")
	if err != nil {
		return err
	}
	o.client = ollamaApi.NewClient(u, http.DefaultClient)
	// g.client, err = api.ClientFromEnvironment()
	return nil
}

func (o *Ollama) GenerateCommit(ctx context.Context, diff string) (string, error) {
	log.Println("Requesting...")
	stream := false
	req := &ollamaApi.GenerateRequest{
		Model:  "qwen2.5-coder:7b",
		Prompt: util.GeneralPrompt + (diff),
		Options: map[string]any{
			"temperature": 0,
		},

		// set streaming to false
		Stream: &stream,
	}
	var result string
	respFunc := func(resp ollamaApi.GenerateResponse) error {
		result = resp.Response
		return nil
	}

	err := o.client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", err
	}
	return result, nil
}
