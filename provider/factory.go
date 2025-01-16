package provider

import (
	"fmt"
	"tdd-go/config"
	"tdd-go/llm"
)

func NewTextGenerator(providerName string, config config.Config) (llm.TextGenerator, error) {
	switch providerName {
	case "oai":
		return llm.OAI{
			BaseURL: config.OpenAI.BaseURL,
			Key:     config.OpenAI.Key,
		}, nil
	case "anthropic":
		return llm.Claude{
			BaseURL: config.Anthropic.BaseURL,
			Key:     config.Anthropic.Key,
		}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", providerName)
	}
}
