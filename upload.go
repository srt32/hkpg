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
	var bucketName = os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("S3_BUCKET_NAME must be set")
	}

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
	sesh, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create AWS session: %v", err)
	}
	svc := s3.New(sesh, cfg)

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("%v", err)
	}

	buffer := new(bytes.Buffer)
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatalf("failed to seek file: %v", err)
	}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		log.Fatalf("file read failed: %v", err)
	}
	fileBytes := bytes.NewReader(buffer.Bytes())
	fileType := http.DetectContentType(buffer.Bytes())

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(file.Name()),
		Body:          fileBytes,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(fileType),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		log.Fatalf("bad response: %s", err)
	}

	var etag = awsutil.StringValue(resp)
	return etag, nil
}
