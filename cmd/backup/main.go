package main

import (
	"log"
	"os"

	"github.com/srt32/hkpg"
	"github.com/srt32/hkpg/download"
	"github.com/srt32/hkpg/heroku"
)

func main() {
	var herokuAppName = os.Getenv("HEROKU_APP_NAME")
	if herokuAppName == "" {
		log.Fatalf("HEROKU_APP_NAME must be set")
	}

	var newestTransfer = heroku.GetTransfers(herokuAppName)
	var publicUrl = heroku.GetPublicUrl(newestTransfer, herokuAppName)
	log.Printf("Success! %v", publicUrl)

	file, err := download.DownloadUrl(publicUrl.Url)
	if err != nil {
		log.Fatalf("download failed, %v", err)
	}
	log.Printf("Success! %v", file)

	uploadedFileUrl, err := hkpg.Upload(file, "upload-file-name-foo")
	if err != nil {
		log.Fatalf("upload failed, %v", err)
	}
	log.Printf("Success! %v", uploadedFileUrl)

	os.Exit(0)
}
