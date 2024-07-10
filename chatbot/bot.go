package chatbot

import (
	"log"
	"os"

	"github.com/ayushthe1/streak/database"
	"github.com/ayushthe1/streak/models"
	"github.com/joho/godotenv"
)

var CredentialFile string
var ProjectID string
var BotPassword string

const ChatbotUsername = "ChatBot"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Couldn't load env file")
	}

	credentialFile := os.Getenv("CREDENTIAL_FILE_PATH")
	if credentialFile == "" {
		log.Fatalf("Empty credenrtial file provided")
	}

	projectID := os.Getenv("DIALOGFLOW_PROJECT_ID")
	if projectID == "" {
		log.Fatalf("No projectID provided")
	}

	botPassword := os.Getenv("BOT_PASSWORD")
	BotPassword = botPassword

	CredentialFile = credentialFile
	ProjectID = projectID
}

// function to create 'chatbot' as a user. It checks DB to create 'chatbot' if it dont exists,
func CreateChatBotUser() error {
	var botUser models.User

	database.DB.Where("username=?", ChatbotUsername).First(&botUser)
	if botUser.Id != 0 {
		log.Println("Bot already exists, Not creating it.")
		return nil
	}

	bot := models.User{
		Username: ChatbotUsername,
	}
	bot.SetPassword(BotPassword)
	result := database.DB.Create(&bot)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	log.Println("Bot User successfully created")
	return nil
}
