package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ayushthe1/streak/database"
	"github.com/ayushthe1/streak/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var ErrNoChatHistory = "No chat history found"

// check whether a given username exists or not
func VerifyContactHandler(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		log.Println("Unable to parse body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	data, err := GetUserByUsername(user.Username)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid Username",
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "username is valid",
		"data":    data,
	})
}

// Handler function to return chat history between 2 users
func ChatHistoryHandler(c *fiber.Ctx) error {

	u1 := c.Query("u1")
	u2 := c.Query("u2")

	_, err := GetUserByUsername(u1)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid Username",
		})
	}

	_, err = GetUserByUsername(u2)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid Username",
		})
	}

	allChats, err := GetAllChats(u1, u2)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Status(200)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("No chat history found between %s and %s", u1, u2),
			})
		}
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": fmt.Errorf("error : %v", err.Error()),
		})
	}

	if u2 == "ChatBot" {

		chatBotIntroText := `Hello! I'm Streak's AI assistant, ready to enhance your chat experience. Here's how I can help you:

- Get real-time weather updates for any city
- Answer general knowledge questions on a wide range of topics
- Provide quick explanations and definitions
- Offer suggestions and recommendations

I'm constantly learning and evolving. Stay tuned for exciting new capabilities coming soon!

How can I assist you today?`

		msg := models.Chat{
			From:      u2,
			To:        u1,
			Msg:       chatBotIntroText,
			Timestamp: time.Now().Unix(),
		}
		allChats = append(allChats, msg)
	}

	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Chat history found",
		"chats":   allChats,
		"total":   len(allChats),
	})

}

// Returns all the contacts of a user (for now it returns all the users)
func ContactHandler(c *fiber.Ctx) error {
	username := c.Query("username")
	log.Println(username)

	allUsers, err := GetContacts()
	if err != nil {
		log.Println("Error 3")
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Couldn't get all contacts",
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"Status": true,
		"data":   allUsers,
		"Total":  len(allUsers),
	})
}

// Return the last 10 aactivities
func ActivityHandler(c *fiber.Ctx) error {
	activities, err := GetActivities()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Couldn't get the activities",
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"activities": activities,
	})
}

// function to get the latest 10 actiities
func GetActivities() (activities []models.ActivityEvent, err error) {

	var results []models.ActivityEvent
	err = database.DB.Order("timestamp desc").Limit(10).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetUserByUsername(username string) (*models.User, error) {

	var user models.User
	result := database.DB.Where("username=?", username).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("username: %v not found", username)
		}
		return nil, result.Error
	}

	return &user, nil
}

// Get all the chat history between 2 users. Exepects the usernames provided to already be verified
func GetAllChats(u1 string, u2 string) ([]models.Chat, error) {

	var allChats []models.Chat

	// result := database.DB.Where("from=?", username1).Or("username=?", username2).Find(&allChats)
	result := database.DB.Where(`("from" = ? AND "to" = ?) OR ("from" = ? AND "to" = ?)`, u1, u2, u2, u1).Find(&allChats)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf(ErrNoChatHistory)
		}
		return nil, result.Error
	}

	return allChats, nil

}

func GetContacts() ([]models.User, error) {

	var allUsers []models.User

	result := database.DB.Select("username", "last_activity").Find(&allUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return allUsers, nil

}

// func CreateChat(chatMsg *models.Chat) (interface{}, error) {

// 	result := database.DB.Create(chatMsg)
// 	if result.Error != nil {
// 		return "", result.Error
// 	}
// 	log.Printf("message : %s saved in DB : ", chatMsg.Msg)

// 	return chatMsg.Id, nil
// }
