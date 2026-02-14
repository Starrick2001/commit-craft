package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

	if g.client != nil {
		return nil
	}
	g.client, err = genai.NewClient(ctx, &genai.ClientConfig{APIKey: g.Config.APIKey})
	if err != nil {
		return err
	}

	return nil
}

func (g *GeminiAdapter) GenerateCommit(ctx context.Context, diff string) (*Output, error) {
	var output *Output
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

func (g *GeminiAdapter) GetListModel(ctx context.Context) ([]*config.ModelOption, error) {
	if err := g.InitClient(ctx); err != nil {
		return nil, errors.New("can not init llm client with message: " + err.Error())
	}

	geminiModels, err := g.client.Models.List(ctx, &genai.ListModelsConfig{PageSize: 10})
	if err != nil {
		return nil, err
	}

	modelOptions := []*config.ModelOption{
		{Name: fmt.Sprintf("Default Gemini Model (%v)", util.GEMINI_DEFAULT_MODEL), Code: util.GEMINI_DEFAULT_MODEL},
	}

	for _, geminiModel := range geminiModels.Items {
		for _, action := range geminiModel.SupportedActions {
			if action == "generateContent" {
				modelOptions = append(modelOptions, &config.ModelOption{Name: geminiModel.DisplayName, Code: geminiModel.Name, Description: geminiModel.Description})
			}
		}
	}
	if len(modelOptions) == 0 {
		return nil, errors.New("can not find any models")
	}
	return modelOptions, nil
}
