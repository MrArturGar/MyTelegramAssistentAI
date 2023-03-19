package messageLayer

import (
	"MyTelegramAssistentAI/src/models"
	"MyTelegramAssistentAI/src/modules"
	"MyTelegramAssistentAI/src/services/dataRepository"
	"MyTelegramAssistentAI/src/services/gptservice"
	"log"
)

var repository *dataRepository.DataRepository
var gptclient *gptservice.ChatGPT
var modulePool *modules.Pool

func Start() {
	repository = dataRepository.GetInstance()
	gptclient = gptservice.GetInstance()
	modulePool = modules.GetInstance()
}

func SendReuest(chatid int64, request *string) (response *string, attachment *string) {
	if isSpecRequest(request) {
		return SpecCommandReplacer(request)
	} else {
		repository.AddMessage(chatid, models.User, *request)
		conversation := repository.GetConversation(chatid)
		response, err := gptclient.Send(conversation.Messages)
		log.Println("GPT reponse: " + *response)

		if err != nil {
			log.Fatal(err)
		} else {
			repository.AddMessage(chatid, models.Assistent, *response)
		}

		return SpecCommandReplacer(response)
	}

}

func isSpecRequest(command *string) bool {
	return (*command)[0:1] == "/"
}

func executeModule(param *string, value *string) *string {
	module, err := modulePool.Get(*param)

	if err != nil {
		log.Print(err.Error())
		return param //TODO
	}

	return module.Execute(value)
}

func getSpecCommand(request *string) (param *string, value *string) {
	for i := 0; i < len(*request); i++ {
		log.Print((*request)[i])
		if (*request)[i] == ' ' {
			param := (*request)[1:i]
			value := (*request)[i+1 : len(*request)]

			return &param, &value
		}
	}
	paramSimple := (*request)[1:]
	return &paramSimple, nil
}

func SpecCommandReplacer(response *string) (*string, *string) {
	startSpecCmdPos := -1
	var responseBuffer, param, value string
	for i := 0; i < len(*response); i++ {
		char := (*response)[i]

		if char == '/' {
			startSpecCmdPos = i
			if i != 0 {
				responseBuffer += (*response)[0 : i-1]
			}
		}

		if startSpecCmdPos != -1 {
			if char == ' ' || char == '.' || char == '\n' {
				if param == "" {
					param = (*response)[startSpecCmdPos+1 : i]
					startSpecCmdPos = i + 1
				}

				if char != ' ' {
					value = (*response)[startSpecCmdPos:i]
					responseBuffer += (*response)[i:len(*response)]
				}
			}

			if (i+1) == len(*response) && value == "" {
				value = (*response)[startSpecCmdPos : i+1]
				if i != len(*response)-1 {
					responseBuffer += (*response)[i : len(*response)-1]
				}
			}
		}

	}

	if param == "" {
		return response, &param
	}
	return &responseBuffer, executeModule(&param, &value)

}
