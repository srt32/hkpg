package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// example:
// {"uuid":"","num":838,"from_name":"DATABASE","from_type":"pg_dump","from_url":"","to_name":"SCHEDULED BACKUP","to_type":"gof3r","to_url":"","options":{},"source_bytes":41722568876,"processed_bytes":6954926696,"succeeded":true,"warnings":0,"created_at":"2017-08-31 10:28:08 +0000","started_at":"2017-08-31 10:29:21 +0000","canceled_at":null,"updated_at":"2017-08-31 10:47:04 +0000","finished_at":"2017-08-31 10:46:59 +0000","deleted_at":null,"purged_at":null,"num_keep":5,"schedule":{"uuid":""}}

type transfers_struct struct {
	Test string
}

// example:
// {"expires_at":"2017-08-31 23:51:02 +0000","url":""}

type transfer_struct struct {
}

func main() {
	var herokuAuthToken = os.Getenv("HEROKU_AUTH_TOKEN")
	if herokuAuthToken == "" {
		log.Fatalf("HEROKU_AUTH_TOKEN must be set")
	}

	var herokuAppName = os.Getenv("HEROKU_APP_NAME")
	if herokuAppName == "" {
		log.Fatalf("HEROKU_APP_NAME must be set")
	}

	var transfersUrl = fmt.Sprintf(
		"https://api.heroku.com/client/v11/apps/%s/transfers",
		herokuAppName,
	)

	// TODO / NEXT: add -H "Authorization: Bearer $TUTORIAL_KEY"
	// TODO / NEXT: https://stackoverflow.com/a/32540356/1949363
	resp, err := http.Post(
		transfersUrl,
		"application/vnd.heroku+json; version=3",
		bytes.NewBuffer([]byte{}),
	)

	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	var t transfers_struct
	err = decoder.Decode(&t)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf(fmt.Sprintf("transfers call failed with %v", resp.StatusCode))
	}
	log.Println(t)

	// get entry with highest `num`
	// 2. POST /client/v11/apps/backerkit-staging2/transfers/{highest_num}/actions/public-url
	// pluck out `url` field

	// - download it
	// - upload it to s3
	// https://github.com/aws/aws-sdk-go/tree/master/service/s3

	// report success
	// call honeybadger (optionally)
	log.Printf("Success!")
}
