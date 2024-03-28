package birthday

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Generator struct {
	client *openai.Client
}

type Options struct {
	OpenAIKey string
}

const maxTokens = 500
const sytemPrompt = `
Твоя задача оригинально и смешно поздравить человека с днем рождения, с учетом полученных данных о нем.
Представь что ты работаешь с этим человеком уже давно и у тебя сегодня хорошее настроение`

func NewGenerator(options Options) *Generator {
	return &Generator{
		client: openai.NewClient(options.OpenAIKey),
	}
}

func (g *Generator) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			MaxTokens: maxTokens,
			Model:     openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: sytemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to ask OpenAI: %w", err)
	}

	answer := resp.Choices[0].Message.Content
	if answer == "" {
		return "", fmt.Errorf("OpenAI returned empty answer")
	}

	return answer, nil
}

func (g *Generator) Prompt(name, description string) string {
	return name + " " + description
}
