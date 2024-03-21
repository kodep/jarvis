package birthday

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type UserData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	NickName    string `json:"nick_name"`
}

const MaxTokens = 500
const SystemPrompt = `
Твоя задача оригинально и смешно поздравить человека с днем рождения, с учетом полученных данных о нем`

func NewClient() (*openai.Client, error) {
	key := os.Getenv("OPEN_AI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("%s environment variable is missing", "OPEN_AI_API_KEY")
	}

	client := openai.NewClient(key)

	return client, nil
}

func GetMessage(r *http.Request) (string, error) {
	userData := &UserData{}

	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body: %w", err)
	}

	err = json.Unmarshal(byteBody, userData)

	if err != nil {
		return "", fmt.Errorf("cannot parse body: %w", err)
	}

	prompt := userData.Name + " " + userData.Description

	gptMessage, err := generateCongratulation(prompt)
	if err != nil {
		return "", err
	}

	message := gptMessage

	if userData.NickName != "" {
		message = "@" + userData.NickName + " " + message
	}

	return message, nil
}

func generateCongratulation(userPrompt string) (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", fmt.Errorf("cannot get open ai client: %w", err)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			MaxTokens: MaxTokens,
			Model:     openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: SystemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	answer := resp.Choices[0].Message.Content

	if answer == "" {
		return "", errors.New("empty generation")
	}

	return answer, nil
}
