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

type Config struct {
	Provider     string
	Model        string
	APIKey       string
	PrefixCommit string
	// For Gemini
	ThinkingBudget int32
}

func BuildConfig() (*Config, error) {
	apiKey, isEnvAPIKeyFound := os.LookupEnv("COMMIT_CRAFT_GEMINI_KEY")
	if !isEnvAPIKeyFound {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Can not find Gemini Api key in your env, what is your Gemini Api Key?").Value(&apiKey),
			),
		)
		if err := form.Run(); err != nil {
			log.Println(`Failed to build "huh" form ` + err.Error())
			return nil, err
		}
	}
	config := &Config{APIKey: apiKey, PrefixCommit: os.Getenv("COMMIT_CRAFT_PREFIX_COMMIT")}
	providerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a provider:").
				Options(
					huh.NewOption("Gemini", (GeminiClient)),
					huh.NewOption("Ollama", (OllamaClient)),
				).
				Value(&config.Provider)))
	err := providerForm.Run()
	if err != nil {
		return nil, err
	}
	switch config.Provider {
	case GeminiClient:
		modelForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose a model:").
					Options(
						huh.NewOption("gemini-2.0-flash", "gemini-2.0-flash"),
						huh.NewOption("gemini-2.5-flash", "gemini-2.5-flash")).
					Value(&config.Model)),
		)
		if err := modelForm.Run(); err != nil {
			log.Println(`Failed to build "huh" form ` + err.Error())
			return nil, err
		}
	}
	return config, nil
}
