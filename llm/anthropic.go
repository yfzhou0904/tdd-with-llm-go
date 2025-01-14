package llm

import (
	"context"
	"log/slog"

	anthropic "github.com/anthropics/anthropic-sdk-go" // imported as anthropic
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Claude struct {
	BaseURL string
	Key     string
}

func (c Claude) GenerateText(prompt string) (string, error) {
	client := anthropic.NewClient(
		option.WithAPIKey(c.Key),
		option.WithBaseURL(c.BaseURL),
	)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(2048)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		panic(err.Error())
	}

	response := ""
	for _, contentBlock := range message.Content {
		if contentBlock.Type != "text" {
			slog.Warn("unexpected content block type", "contentBlock", contentBlock)
		}
		response += contentBlock.Text
	}
	return response, nil
}
