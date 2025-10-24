package provider

import (
	"context"
	"errors"
	"fmt"
	"log"

	"commit-craft/config"

	"github.com/charmbracelet/huh/spinner"
)

type Provider interface {
	InitClient(ctx context.Context) error
	GenerateCommit(ctx context.Context, diff string) (*Output, error)
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
	log.Println("GetClient called with provider:", cfg.Provider)
	llmOption := cfg.Provider
	var provider Provider
	switch llmOption {
	case config.OllamaClient:
		log.Println("Selected Ollama provider")
		provider = &OllamaAdapter{Config: cfg}
	case config.GeminiClient:
		log.Println("Selected Gemini provider")
		provider = &GeminiAdapter{Config: cfg}
	default:
		log.Printf("Unknown or missing provider: '%s', defaulting to Gemini for safety or erroring", llmOption)
		// For now, let's return error to avoid panic
		return nil, errors.New("unknown provider: " + llmOption)
	}

	return &Client{provider: provider, cfg: cfg}, nil
}

// Combine the prefix, title, and description into a complete commit message.
func (c *Client) buildCommit(output *Output) string {
	log.Println("Building commit message")
	return fmt.Sprintf("%s%s\n\n%s", c.cfg.PrefixCommit, output.Title, output.Description)
}

func (c *Client) GenerateCommit(ctx context.Context, diff string) (string, error) {
	log.Println("GenerateCommit called")
	if c.provider == nil {
		return "", errors.New("provider is nil")
	}
	if err := c.provider.InitClient(ctx); err != nil {
		log.Println("InitClient failed:", err)
		return "", errors.New("can not init llm client with message: " + err.Error())
	}

	var result *Output
	log.Println("Requesting from provider...")
	err := spinner.New().Title("Requesting...").Context(ctx).ActionWithErr(func(context.Context) error {
		var err error
		result, err = c.provider.GenerateCommit(ctx, diff)
		return err
	}).Accessible(false).Run()
	if err != nil {
		log.Println("Provider GenerateCommit failed:", err)
		return "", errors.New("can not generate commit with message: " + err.Error())
	}
	log.Println("Commit generated successfully")

	return c.buildCommit(result), nil
}
