package hkpg

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TODO: drop the fileName?
func Upload(file *os.File, fileName string) (string, error) {
	creds := credentials.NewEnvCredentials{}
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("bad credentials: %s", err)
	}

	// TODO: accept region as arg
	cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(creds)

	svc := s3.New(session.New(), cfg)

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	path := "/backups/" + file.Name()
	params := &s3.PutObjectInput{
		Bucket:        aws.String("nameofBucketHere"), // TODO: make configurable
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		log.Fatalf("bad response: %s", err)
	}

	log.Printf("response %s", awsutil.StringValue(resp))

	return awsutil.StringValue(resp), nil
}
