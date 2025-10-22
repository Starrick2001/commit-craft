package provider

import (
	"context"
	"encoding/json"

	"commit-craft/config"
	"commit-craft/util"

	"google.golang.org/genai"
)

type GeminiAdapter struct {
	Config *config.Config
	client *genai.Client
}

func (g *GeminiAdapter) InitClient(ctx context.Context) error {
	var err error
	g.client, err = genai.NewClient(ctx, &genai.ClientConfig{APIKey: g.Config.APIKey})
	if err != nil {
		return err
	}
	return nil
}

func (g *GeminiAdapter) GenerateCommit(ctx context.Context, diff string) (*LLMResponse, error) {
	var output *LLMResponse
	result, err := g.client.Models.GenerateContent(
		ctx,
		g.Config.Model,
		genai.Text(util.GeneralPrompt+(diff)),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &g.Config.ThinkingBudget, // Disables thinking
			},
			ResponseMIMEType: "application/json",
			ResponseSchema:   &genai.Schema{Type: genai.TypeObject, Properties: map[string]*genai.Schema{"title": {Type: genai.TypeString}, "description": {Type: genai.TypeString}}},
		},
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result.Text()), &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
