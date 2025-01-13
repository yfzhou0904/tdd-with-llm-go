package llm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OAI struct {
	BaseURL string
	Key     string
}

func (o OAI) GenerateText(prompt string) (string, error) {
	client := openai.NewClient(
		option.WithBaseURL(o.BaseURL),
		option.WithAPIKey(o.Key),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelChatgpt4oLatest),
	})
	if err != nil {
		return "", err
	} else if len(chatCompletion.Choices) == 0 {
		slog.Error("no choices returned", "chatCompletion", chatCompletion)
		return "", fmt.Errorf("no choices returned")
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
