package provider

import (
	"context"
	"errors"
	"fmt"

	"commit-craft/config"

	"github.com/charmbracelet/huh/spinner"
)

type Provider interface {
	InitClient(ctx context.Context) error
	GenerateCommit(ctx context.Context, diff string) (*Output, error)
	GetListModel(ctx context.Context) ([]*config.ModelOption, error)
}

type Client struct {
	provider Provider
	cfg      *config.Config
}

type Output struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetClient(cfg *config.Config) (*Client, error) {
	llmOption := cfg.Provider
	var provider Provider
	switch llmOption {
	case config.OllamaClient:
		provider = &OllamaAdapter{Config: cfg}
	case config.GeminiClient:
		provider = &GeminiAdapter{Config: cfg}
	}

	return &Client{provider: provider, cfg: cfg}, nil
}

// Combine the prefix, title, and description into a complete commit message.
func (c *Client) buildCommit(output *Output) string {
	return fmt.Sprintf("%s%s\n\n%s", c.cfg.PrefixCommit, output.Title, output.Description)
}

func (c *Client) GenerateCommit(ctx context.Context, diff string) (string, error) {
	if err := c.provider.InitClient(ctx); err != nil {
		return "", errors.New("can not init llm client with message: " + err.Error())
	}

	var result *Output
	err := spinner.New().Title("Requesting...").Context(ctx).ActionWithErr(func(context.Context) error {
		var err error
		result, err = c.provider.GenerateCommit(ctx, diff)
		return err
	}).Accessible(false).Run()
	if err != nil {
		return "", errors.New("can not generate commit with message: " + err.Error())
	}

	return c.buildCommit(result), nil
}

func (c *Client) GetListModel(ctx context.Context) ([]*config.ModelOption, error) {
	return c.provider.GetListModel(ctx)
}
