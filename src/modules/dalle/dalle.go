package dalle

import (
	"MyTelegramAssistentAI/src/config"
	"context"
	"log"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

type dalle struct {
	client *openai.Client
}

var instance *dalle
var once sync.Once

func GetInstance() *dalle {
	once.Do(func() {
		instance = &dalle{}
		open(instance)
	})
	return instance
}

func open(gpt *dalle) {
	gpt.client = openai.NewClient(config.GetValue("OPENAI_TOKEN"))
}

func (dalle *dalle) Send(prompt string) *string {
	resp, err := dalle.client.CreateImage(context.Background(),
		openai.ImageRequest{
			Prompt: prompt,
			N:      1,
			Size:   "1024x1024",
		},
	)

	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return &resp.Data[0].URL
}

func (dalle *dalle) Execute(request *string) *string {
	return dalle.Send(*request)
}
