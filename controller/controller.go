package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"

	"commandgpt/model/resource"
)

var ErrResponseParseFailed = errors.New("failed to parse response")

type Controller struct {
	apiKey string
	clt    *openai.Client
}

func New() Controller {
	return Controller{
		apiKey: os.Getenv("OPENAI_API_KEY"),
	}
}

func (c *Controller) Init() error {
	if c.apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is not set")
	}

	c.clt = openai.NewClient(c.apiKey)

	return nil
}

func (c *Controller) buildMessages(query string, number int) []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, len(prompts)+1)

	copy(messages, prompts)

	messageText, err := json.Marshal(resource.Query{
		Query:  query,
		Number: number,
	})
	if err != nil {
		panic(err)
	}

	messages[len(prompts)] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: string(messageText),
	}

	return messages
}

func (c *Controller) GetSuggestion(model, query string, number int) (*resource.Advices, error) {
	messages := c.buildMessages(query, number)

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	resp, err := c.clt.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var advices resource.Advices
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &advices); err != nil {
		return nil, ErrResponseParseFailed
	}

	return &advices, nil
}

func (c *Controller) TestGetSuggestion(_, _ string, _ int) (*resource.Advices, error) {
	time.Sleep(5 * time.Second)

	advices := resource.Advices{
		Advices: []resource.Advice{
			{
				Command:     "cli 1",
				Description: "use cli",
			},
			{
				Command:     "cli2",
				Description: "cli2 description",
			},
			{
				Command:     "cli3",
				Description: "cli3 description",
			},
			{
				Command:     "cli4",
				Description: "cli4 description",
			},
			{
				Command:     "cli5",
				Description: "cli5 description",
			},
		},
	}

	return &advices, nil
}
