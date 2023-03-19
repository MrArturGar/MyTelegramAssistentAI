package gptservice

import (
	"MyTelegramAssistentAI/src/config"
	"context"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	client *openai.Client
}

var instance *ChatGPT
var once sync.Once

func GetInstance() *ChatGPT {
	once.Do(func() {
		instance = &ChatGPT{}
		open(instance)
	})
	return instance
}

func open(gpt *ChatGPT) {
	gpt.client = openai.NewClient(config.GetValue("OPENAI_TOKEN"))
}

func (gpt *ChatGPT) Send(conversation []openai.ChatCompletionMessage) (*string, error) {

	resp, err := gpt.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: conversation,
		},
	)

	if err != nil {
		return nil, err
	}
	return &(resp.Choices[0].Message.Content), nil
}
