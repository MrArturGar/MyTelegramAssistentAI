package models

import openai "github.com/sashabaranov/go-openai"

type Conversation struct {
	UserId   int64
	Messages []openai.ChatCompletionMessage
}

type Role string

const (
	User      Role = "user"
	Assistent Role = "assistant"
	System    Role = "system"
)
