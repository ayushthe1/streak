package s3

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3session *session.Session

func init() {
	ConnectToS3()
}

func ConnectToS3() {
	s3Session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))

	S3session = s3Session
	log.Println("*******CONNECTED TO S3***********")
}

func UploadFileToS3(filepath string, fromUsername string) (s3FileUrl string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open the file %s: %v", filepath, err)
	}

	defer file.Close()

	uploader := s3.New(S3session)
	bucket := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	if len(bucket) == 0 || len(region) == 0 {
		log.Fatalf("No bucket or region provided")
	}

	//name of the file in 3 bucket
	key := fmt.Sprintf("%s%s", fromUsername, file.Name())

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload the file to s3: %v", err)

	}

	// get thr file s3 url
	s3url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
	log.Println("s3Url of file is : ", s3url)

	return s3url, nil
}
