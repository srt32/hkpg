package hkpg

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TODO: drop the fileName?
func Upload(file *os.File, fileName string) (string, error) {
	creds := credentials.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("bad credentials: %s", err)
	}

	// TODO: accept region as arg
	cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(creds)

	svc := s3.New(session.New(), cfg)

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("%v", err)
	}

	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	var bucketName = os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("S3_BUCKET_NAME must be set")
	}
	path := file.Name()

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		log.Fatalf("bad response: %s", err)
	}

	var etag = awsutil.StringValue(resp)
	return etag, nil
}
