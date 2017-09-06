package hkpg

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	sesh := session.Must(session.NewSession(cfg))
	uploader := s3manager.NewUploader(sesh, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024 // 64MB per part
	})

	uploadInput := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    aws.String(file.Name()),
		Body:   file,
	}

	result, err := uploader.Upload(uploadInput)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return result.Location, nil
}
