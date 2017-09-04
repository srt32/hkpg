package main

import (
	"log"
	"os"

	"github.com/srt32/hkpg"
	"github.com/srt32/hkpg/download"
	"github.com/srt32/hkpg/heroku"
)

func main() {
	var herokuAppName = os.Getenv("HKPG_HEROKU_APP_NAME")
	if herokuAppName == "" {
		log.Fatalf("HKPG_HEROKU_APP_NAME must be set")
	}

	var newestTransfer = heroku.GetTransfers(herokuAppName)
	var publicUrl = heroku.GetPublicUrl(newestTransfer, herokuAppName)
	log.Printf("Archive URL fetch successful! %v", publicUrl)

	file, err := download.DownloadUrl(publicUrl.Url, &newestTransfer)
	defer file.Close()
	if err != nil {
		log.Fatalf("download failed, %v", err)
	}
	log.Printf("Download successful! %v", file)

	uploadedETag, err := hkpg.Upload(file)
	if err != nil {
		log.Fatalf("upload failed, %v", err)
	}

	err = os.Remove(file.Name())
	if err != nil {
		log.Fatalf("Upload successful but removal of local backup failed! %v", err)
	}

	log.Printf("Upload successful! %v", uploadedETag)

	os.Exit(0)
}
