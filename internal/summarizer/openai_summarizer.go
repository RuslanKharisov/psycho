package summarizer

import (
	"context"
	"tg-bot/internal/memory"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAISummarizer struct {
	client *openai.Client
}

func New(apiKey string) Summarizer {
	return &OpenAISummarizer{client: openai.NewClient(apiKey)}
}

func (s *OpenAISummarizer) Summarize(userID int64, history []memory.ChatMessage) (string, error) {
	// Формируем prompt для суммаризации
	messages := []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: "Сократи историю диалога, оставив ключевые факты и проблемы пользователя."}}
	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content})
	}
	resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:     "gpt-4o",
		Messages:  messages,
		MaxTokens: 300,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
