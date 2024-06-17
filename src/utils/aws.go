package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func CreateKeyName(userID string) string {
	keyName := fmt.Sprintf("images/%v.webp", userID)
	return keyName
}

func UploadS3FromString(fileName []byte, keyName string, contentType string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(fileName),
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(keyName),
		ContentType: aws.String(contentType),
		Metadata:    map[string]*string{"Content-Disposition": aws.String("attachment")},
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	// req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
	// 	Bucket: aws.String(os.Getenv("AWS_BUCKET")),
	// 	Key:    aws.String(keyName),
	// })
	// url, err := req.Presign(15 * time.Minute)
	// if err != nil {
	// 	return "", err
	// }

	bucket := os.Getenv("AWS_BUCKET")
	region := os.Getenv("AWS_REGION")
	url := "https://" + bucket + ".s3." + region + ".amazonaws.com/" + keyName

	return url, nil
}
