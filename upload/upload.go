package upload

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ayushthe1/streak/channels"
	"github.com/ayushthe1/streak/models"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel
var S3Session *session.Session

type File struct {
	Filepath string
	From     string
	To       string
}

func init() {
	ConnectoRabbitMQ()
	ConnectToS3()
}

func ConnectoRabbitMQ() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	conn, err := amqp.Dial(amqpServerURL)

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	RabbitConn = conn

	channel, err := RabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	RabbitChannel = channel

	log.Println("Connected to RabbitMQ in Upload File")

}

func ConnectToS3() {
	s3Session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	S3Session = s3Session
}

// function to upload file to RabbitMQ
func FileUploadProducer(fileMsg *models.File) error {

	messageBody, err := json.Marshal(fileMsg)
	if err != nil {
		log.Println("error marshalling in fileUploadProducer :", err.Error())
		return err
	}

	err = RabbitChannel.Publish(
		"",
		"file_upload",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		})

	if err != nil {
		return err
	}
	log.Println("Pulished file to RabbitMQ in FileUploadProducer")
	return nil
}

// function to process file upload from RabbitMQ
func FileUploadConsumer() error {

	_, err := RabbitChannel.QueueDeclare(
		"file_upload", // name of the queue
		true,
		false,
		false,
		false,
		nil, //  arguments
	)
	if err != nil {
		log.Println("err :", err.Error())
		return err
	}

	files, err := RabbitChannel.Consume(
		"file_upload",
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Error inside ProcessFileUpload while consuming msg")
		return err
	}

	for f := range files {
		var file models.File
		err := json.Unmarshal(f.Body, &file)
		if err != nil {
			log.Println("Failed to unmarshal chat message:", err)
			continue
		}

		s3FileURL, err := uploadFileToS3(file.TempFilePath, file.From)

		if err != nil {
			log.Printf("error in uploading file to S3 : %v", err)
			return err
		}

		os.Remove(file.TempFilePath)

		channels.Broadcast_S3_FileURL <- s3FileURL
		log.Println("sent to broadcast_s3_fileURL channel")

	}

	return nil

}

func uploadFileToS3(filePath string, fromUsername string) (s3FileURL string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	uploader := s3.New(S3Session)
	bucket := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	if len(bucket) == 0 || len(region) == 0 {
		log.Fatalf("No bucket or region provided")
		return "", fmt.Errorf("No bucket or region")
	}

	// name of file in S3 bucket
	key := fmt.Sprintf("%s/%s", fromUsername, file.Name())

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// get the file S3 url
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
	log.Println("s3URL of file is : ", s3URL)

	return s3URL, nil
}

func SetupFileUpload() {
	go FileUploadConsumer()
}

// func LoadConfig() {

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error load .env file")
// 	}

// 	bucketName := os.Getenv("BUCKET_NAME")
// 	if len(bucketName) == 0 {
// 		log.Fatalf("bucket name not provided")
// 	}

// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		log.Fatalf("failed to load SDK configuration, %v", err)
// 	}

// 	client := s3.NewFromConfig(cfg)
// }
