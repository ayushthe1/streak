package upload

import (
	"log"

	"github.com/ayushthe1/streak/channels"
	"github.com/ayushthe1/streak/kafka"

	"github.com/ayushthe1/streak/models"
)

// function to upload file to kafka
func FileUploadProducer(fileMsg *models.File) (string, error) {

	fileMsg.Type = kafka.FileMsgType
	err := kafka.ProduceEventToKafka(kafka.FileTopic, fileMsg)
	if err != nil {
		log.Println("error uploading the file to kafka  : ", err.Error())
		return "", err
	}

	log.Println("Waiting for msg on fileUrl channel")

	fileUrl := <-channels.Broadcast_S3_FileURL
	log.Println("S3 file url received from channel ", fileUrl)

	return fileUrl, nil
}
