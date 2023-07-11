package controller

import (
	_ "embed"
	"strings"

	"github.com/sashabaranov/go-openai"
)

//go:embed prompt
var promptText string

var prompts = func() []openai.ChatCompletionMessage {
	data := strings.Split(promptText, strings.Repeat("-", 10))

	promptMessages := make([]openai.ChatCompletionMessage, len(data))

	for i := range data {
		data[i] = strings.TrimSpace(data[i])
		split := strings.Split(data[i], "===")
		promptMessages[i] = openai.ChatCompletionMessage{
			Role:    split[0],
			Content: split[1],
		}
	}

	return promptMessages
}()
