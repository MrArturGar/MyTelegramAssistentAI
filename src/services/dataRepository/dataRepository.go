package dataRepository

import (
	"MyTelegramAssistentAI/src/models"
	"database/sql"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"

	_ "github.com/mattn/go-sqlite3"
)

type DataRepository struct {
	dbcontext *sql.DB
}

var instance *DataRepository
var once sync.Once

func GetInstance() *DataRepository {
	once.Do(func() {
		instance = &DataRepository{}
		open(instance)
	})
	return instance
}

func open(db *DataRepository) {
	filename := "dbfile.db"
	dbFound := false

	_, err := os.Stat(filename)
	dbFound = err == nil || os.IsExist(err)

	db.dbcontext, _ = sql.Open("sqlite3", filename)
	if !dbFound {
		_, err := db.dbcontext.Exec("CREATE TABLE Users (id integer, username varchar, dateCreated timestamp ); CREATE TABLE Conversations ( id integer PRIMARY KEY AUTOINCREMENT, userId integer, role varchar, message varchar);")

		if err != nil {
			log.Printf("Failed to create database: %s", err.Error())
			panic(err)
		}
		db.AddUser(0, "bot")                                                                                                                                   //System user
		db.AddMessage(0, openai.ChatMessageRoleSystem, "You are a helpful assistant named Jarvis. Today's date is"+time.Now().Format("2006-01-02 Monday")+".") //Test data for bot
	}
}

func (db *DataRepository) AddMessage(chatid int64, role models.Role, message string) {
	_, err := db.dbcontext.Exec("INSERT INTO Conversations (userId, role, message) VALUES (?, ?, ?)", strconv.FormatInt(chatid, 10), role, message)

	if err != nil {
		log.Printf("Failed to add message [id: %s] to repository: %s", strconv.FormatInt(chatid, 10), err.Error())
		panic(err)
	}
}

func (db *DataRepository) AddUser(chatid int64, username string) {
	_, err := db.dbcontext.Exec("INSERT INTO Users (id, username, dateCreated) VALUES (?, ?, ?)", strconv.FormatInt(chatid, 10), username, strconv.FormatInt(time.Now().UnixNano(), 10))

	if err != nil {
		log.Printf("Failed to add user [name: %s - id: %s] to repository: %s", username, strconv.FormatInt(chatid, 10), err.Error())
		panic(err)
	}
}

func (db *DataRepository) GetConversation(chatid int64) (Conversation models.Conversation) {
	rows, err := db.dbcontext.Query("SELECT role, message FROM Conversations WHERE Conversations.userId = ? OR Conversations.userId = 0", chatid)

	if err != nil {
		log.Printf("Failed to read conversation [id: %s] from repository: %s", strconv.FormatInt(chatid, 10), err.Error())
		panic(err)
	}
	defer rows.Close()

	Conversation = models.Conversation{}
	Conversation.UserId = chatid
	for rows.Next() {
		message := openai.ChatCompletionMessage{}
		rows.Scan(&message.Role, &message.Content)

		Conversation.Messages = append(Conversation.Messages, message)
	}
	return Conversation
}

func (db *DataRepository) GetSystemMessages() (Conversation models.Conversation) {
	rows, err := db.dbcontext.Query("SELECT role, message FROM Conversations WHERE Conversations.userId = 0")

	if err != nil {
		log.Printf("Failed to read system messages from repository: %s", err.Error())
		panic(err)
	}
	defer rows.Close()

	Conversation = models.Conversation{}
	for rows.Next() {
		message := openai.ChatCompletionMessage{}
		rows.Scan(&message.Role, &message.Content)

		Conversation.Messages = append(Conversation.Messages, message)
	}
	return Conversation
}

func (db *DataRepository) DeleteConversation() {
	_, err := db.dbcontext.Exec("DELETE FROM Conversations WHERE Conversations.userId IS NOT 0;")

	if err != nil {
		log.Printf("Failed to delet messages from repository: %s", err.Error())
		panic(err)
	}
}
