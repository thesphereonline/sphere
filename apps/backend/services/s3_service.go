package services

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Initialize AWS session
func initAWSSession() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

// UploadVideo uploads a video file to S3
func UploadVideo(file *multipart.FileHeader) (string, error) {
	s3Client, err := initAWSSession()
	if err != nil {
		return "", err
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Read file contents
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(src); err != nil {
		return "", err
	}

	// Generate S3 object key
	fileKey := fmt.Sprintf("videos/%s", file.Filename)

	// Upload file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:         aws.String(fileKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		log.Println("Failed to upload file:", err)
		return "", err
	}

	// Return public URL
	return fmt.Sprintf("%s/%s", os.Getenv("S3_BUCKET_URL"), fileKey), nil
}
