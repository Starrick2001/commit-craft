package provider

import (
	"context"
	"errors"

	"commit-craft/config"
)

type Provider interface {
	InitClient(ctx context.Context) error
	GenerateCommit(ctx context.Context, diff string) (string, error)
}

type Client struct {
	provider Provider
}

func GetClient(cfg *config.Config) (*Client, error) {
	llmOption := cfg.Provider
	var provider Provider
	switch llmOption {
	case config.OllamaClient:
		provider = &Ollama{Config: cfg}
	case config.GeminiClient:
		provider = &Gemini{Config: cfg}
	}

	return &Client{provider: provider}, nil
}

func (c *Client) GenerateCommit(ctx context.Context, diff string) (string, error) {
	if err := c.provider.InitClient(ctx); err != nil {
		return "", errors.New("can not init llm client with message: " + err.Error())
	}
	result, err := c.provider.GenerateCommit(ctx, diff)
	if err != nil {
		return "", errors.New("can not generate commit with message: " + err.Error())
	}
	return result, nil
}
