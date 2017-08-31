package main

import (
    "log"
    "os"
)

func main () {
  var herokuAuthToken = os.Getenv("HEROKU_AUTH_TOKEN")

  if herokuAuthToken == "" {
    log.Fatalf("HEROKU_AUTH_TOKEN must be set")
  }

  // -  get latest backup url
    // https://github.com/bgentry/heroku-go
    // 0. get token from env
    // curl -X POST https://api.heroku.com/apps \
    // -H "Accept: application/vnd.heroku+json; version=3" \
    // -H "Authorization: Bearer $TUTORIAL_KEY"
    // 1. GET /client/v11/apps/backerkit-staging2/transfers
    // get entry with highest `num`
    // 2. POST /client/v11/apps/backerkit-staging2/transfers/{highest_num}/actions/public-url
    // pluck out `url` field

  // - download it
  // - upload it to s3
    // https://github.com/aws/aws-sdk-go/tree/master/service/s3

  // report success
    // call honeybadger (optionally)
}
