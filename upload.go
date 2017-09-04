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

// Upload accepts a *os.File, uploads it to the specified S3 buckets and
// returns the ETag of the uploaded file. Upload can be called multiple times
// with the same file name.
func Upload(file *os.File) (string, error) {
	creds := credentials.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("bad credentials: %s", err)
	}

	var awsRegion = os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-west-1"
	}

	cfg := aws.NewConfig().WithRegion(awsRegion).WithCredentials(creds)

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
