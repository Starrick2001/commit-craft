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

	if o.client != nil {
		return nil
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
		Model:  o.Config.Model,
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

func (o *OllamaAdapter) GetListModel(ctx context.Context) ([]*config.ModelOption, error) {
	if err := o.InitClient(ctx); err != nil {
		return nil, err
	}
	modelOptions := []*config.ModelOption{}
	models, err := o.client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, model := range models.Models {
		modelOptions = append(modelOptions, &config.ModelOption{Name: model.Name, Code: model.Model, Description: ""})
	}
	return modelOptions, nil
}
