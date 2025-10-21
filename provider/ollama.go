package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"commit-craft/config"
	"commit-craft/util"

	ollamaApi "github.com/ollama/ollama/api"
)

type OllamaAdapter struct {
	Config *config.Config
	client *ollamaApi.Client
}

type FormatProperty struct {
	Type string `json:"type"`
}

type FormatSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]FormatProperty `json:"properties"`
	Required   []string                  `json:"required"`
}

func (o *OllamaAdapter) InitClient(ctx context.Context) error {
	u, err := url.Parse("http://localhost:11434")
	if err != nil {
		return err
	}
	o.client = ollamaApi.NewClient(u, http.DefaultClient)
	// g.client, err = api.ClientFromEnvironment()
	return nil
}

func (o *OllamaAdapter) GenerateCommit(ctx context.Context, diff string) (*Output, error) {
	formatSchema := &FormatSchema{
		Type: "object",
		Properties: map[string]FormatProperty{
			"title":       {Type: "string"},
			"description": {Type: "string"},
		},
		Required: []string{"title", "description"},
	}
	format, err := json.Marshal(formatSchema)
	if err != nil {
		return nil, err
	}
	stream := false
	req := &ollamaApi.GenerateRequest{
		Model:  "qwen2.5-coder:7b",
		Prompt: util.GeneralPrompt + (diff),
		Options: map[string]any{
			"temperature": 0,
		},

		// set streaming to false
		Stream: &stream,
		Format: format,
	}
	var response *Output
	respFunc := func(resp ollamaApi.GenerateResponse) error {
		err = json.Unmarshal([]byte(resp.Response), &response)
		if err != nil {
			return err
		}
		return nil
	}

	err = o.client.Generate(ctx, req, respFunc)
	if err != nil {
		return nil, err
	}
	return response, nil
}
