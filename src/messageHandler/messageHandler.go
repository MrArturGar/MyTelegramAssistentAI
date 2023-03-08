package messagehandler

import (
	"MyTelegramAssistentAI/src/models"
	"MyTelegramAssistentAI/src/services/dataRepository"
	"MyTelegramAssistentAI/src/services/gptservice"
)

var repository *dataRepository.DataRepository
var gptclient *gptservice.ChatGPT

func Start() {
	repository = dataRepository.GetInstance()
	gptclient = gptservice.GetInstance()
}

func SendReuest(chatid int64, request string) string {
	if isSpecCommand(&request) {
		return "isSpecCommand" //TODO
	}

	repository.AddMessage(chatid, models.User, request)
	conversation := repository.GetConversation(chatid)
	response := gptclient.Send(conversation.Messages)
	repository.AddMessage(chatid, models.Assistent, response)
	return response
}

func isSpecCommand(command *string) bool {
	return (*command)[0:1] == "/"
}
