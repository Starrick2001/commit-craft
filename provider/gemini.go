package provider

import (
	"context"

	"commit-craft/config"
	"commit-craft/util"

	"google.golang.org/genai"
)

type Gemini struct {
	Config *config.Config
	client *genai.Client
}

func (g *Gemini) InitClient(ctx context.Context) error {
	var err error
	g.client, err = genai.NewClient(ctx, &genai.ClientConfig{APIKey: g.Config.APIKey})
	if err != nil {
		return err
	}
	return nil
}

func (g *Gemini) GenerateCommit(ctx context.Context, diff string) (string, error) {
	result, err := g.client.Models.GenerateContent(
		ctx,
		g.Config.Model,
		genai.Text(util.GeneralPrompt+(diff)),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &g.Config.ThinkingBudget, // Disables thinking
			},
		},
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
