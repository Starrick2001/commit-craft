package config

import (
	"log"
	"os"

	"github.com/charmbracelet/huh"
)

const (
	GeminiClient = "gemini"
	OllamaClient = "ollama"
)

type ModelOption struct {
	Name        string
	Code        string
	Description string
}

type Config struct {
	Provider     string
	Model        string
	APIKey       string
	PrefixCommit string
	// For Gemini
	ThinkingBudget int32
}

func BuildConfig() (*Config, error) {
	config := &Config{PrefixCommit: os.Getenv("COMMIT_CRAFT_PREFIX_COMMIT"), ThinkingBudget: 0}

	if err := config.GetAPIKey(); err != nil {
		return nil, err
	}

	if err := config.ChooseProvider(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) GetAPIKey() error {
	apiKey, isEnvAPIKeyFound := os.LookupEnv("COMMIT_CRAFT_GEMINI_KEY")
	if !isEnvAPIKeyFound {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Can not find Gemini Api key in your env, what is your Gemini Api Key?").Value(&apiKey),
			),
		)
		if err := form.Run(); err != nil {
			log.Println(`Failed to build "huh" form ` + err.Error())
			return err
		}
	}

	c.APIKey = apiKey
	return nil
}

func (c *Config) ChooseProvider() error {
	providerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a provider:").
				Options(
					huh.NewOption("Gemini", (GeminiClient)),
					huh.NewOption("Ollama", (OllamaClient)),
				).
				Value(&c.Provider)))
	err := providerForm.Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) ChooseModel(modelOptions []*ModelOption) error {
	formOptions := []huh.Option[string]{}

	for _, modelOption := range modelOptions {
		formOptions = append(formOptions, huh.Option[string]{Value: modelOption.Code, Key: modelOption.Name})
	}

	modelForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a model:").
				Options(
					formOptions...).
				Value(&c.Model)),
	)

	if err := modelForm.Run(); err != nil {
		log.Println(`Failed to build "huh" form: ` + err.Error())
		return err
	}

	return nil
}
