package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	// Create a temporary file within our temp-images directory
	tempDir := "/tmp/streakfile" // replace with actual path
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create temp directory")
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	log.Printf("The temp file path is %s", tempFilePath)

	// Save the file to the temporary path
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
	}

	fileMsg := models.File{
		TempFilePath: tempFilePath,
		From:         senderUsername,
		To:           receiverUsername,
	}

	// //TODO: send some context here
	s3FileUrl, err := upload.FileUploadProducer(&fileMsg)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to upload the file to s3",
			"error":   err.Error(),
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "file uploaded to S3",
		"fileUrl": s3FileUrl,
	})

}
