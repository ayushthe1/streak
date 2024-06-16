package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ayushthe1/streak/channels"
	"github.com/ayushthe1/streak/models"
	"github.com/ayushthe1/streak/upload"
	"github.com/gofiber/fiber/v2"
)

func FileUploadHandler(c *fiber.Ctx) error {

	senderUsername := c.FormValue("sender")
	receiverUsername := c.FormValue("receiver")

	if senderUsername == "" || receiverUsername == "" {
		return c.Status(http.StatusBadRequest).SendString("Sender and Receiver are required")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to get file")
	}

	tempFilePath := fmt.Sprintf("/tmp/streak/%s", file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	fileMsg := models.File{
		TempFilePath: tempFilePath,
		From:         senderUsername,
		To:           receiverUsername,
	}

	//TODO: send some context here
	err = upload.FileUploadProducer(&fileMsg)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"unable to upload file to rabbitMQ": err.Error(),
		})
	}

	log.Println("Waiting for s3_URL on channel in FileUploadHandler")
	s3_URL := <-channels.Broadcast_S3_FileURL
	log.Println("S3_URL received on channel in FileUploadHandler :", s3_URL)

	fileMsg.S3_File_URL = s3_URL

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "file uploaded to S3",
		"fileMsg": fileMsg,
	})

}
