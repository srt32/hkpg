package main

import (
	"log"
	"os"

	"github.com/srt32/hkpg"
)

func main() {
	var herokuAppName = os.Getenv("HEROKU_APP_NAME")
	if herokuAppName == "" {
		log.Fatalf("HEROKU_APP_NAME must be set")
	}

	var newestTransfer = heroku.GetTransfers(herokuAppName)
	log.Printf("Success! %v", newestTransfer)
	os.Exit(0)

	// 2. POST /client/v11/apps/{herokuAppName}/transfers/{highest_num}/actions/public-url
	// pluck out `url` field

	// - download it
	// - upload it to s3
	// https://github.com/aws/aws-sdk-go/tree/master/service/s3
}
