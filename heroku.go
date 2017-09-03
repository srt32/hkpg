package heroku

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
)

// example:
// [{
//   "uuid":""
//   "num":838,
//   "from_name":"DATABASE",
//   "from_type":"pg_dump",
//   "from_url":"",
//   "to_name":"SCHEDULED BACKUP",
//   "to_type":"gof3r",
//   "to_url":"",
//   "options":{},
//   "source_bytes":41722568876,
//   "processed_bytes":6954926696,
//   "succeeded":true,
//   "warnings":0,
//   "created_at":"2017-08-31 10:28:08 +0000",
//   "started_at":"2017-08-31 10:29:21 +0000",
//   "canceled_at":null,
//   "updated_at":"2017-08-31 10:47:04 +0000",
//   "finished_at":"2017-08-31 10:46:59 +0000",
//   "deleted_at":null,
//   "purged_at":null,
//   "num_keep":5,
//   "schedule":{"uuid":""}
// }]

type Transfer struct {
	FinishedAt string `json:"finished_at"`
	FromName   string `json:"from_name"`
	FromType   string `json:"from_type"`
	Num        int
	Succeeded  bool
	ToType     string `json:"to_type"`
	CreatedAt  string `json:"created_at"`
}

type TransfersList []Transfer

func (t TransfersList) Len() int {
	return len(t)
}
func (t TransfersList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TransfersList) Less(i, j int) bool {
	return t[i].CreatedAt < t[j].CreatedAt
}

type PublicUrl struct {
	ExpiresAt string `json:"expires_at"`
	Url       string
}

// example: PublicUrl
// {"expires_at":"2017-08-31 23:51:02 +0000","url":""}

func GetTransfers(appName string) Transfer {
	var herokuAuthToken = os.Getenv("HEROKU_AUTH_TOKEN")
	if herokuAuthToken == "" {
		log.Fatalf("HEROKU_AUTH_TOKEN must be set")
	}

	var transfersUrl = fmt.Sprintf(
		"https://postgres-api.heroku.com/client/v11/apps/%s/transfers",
		appName,
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", transfersUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("", herokuAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	var transfers = make([]Transfer, 0)
	err = decoder.Decode(&transfers)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf(fmt.Sprintf("transfers call failed with %v", resp.StatusCode))
	}

	completedValidTransfers := make([]Transfer, 0)
	for _, transfer := range transfers {
		if transfer.Succeeded == true && transfer.ToType == "gof3r" {
			completedValidTransfers = append(completedValidTransfers, transfer)
		}
	}
	if len(completedValidTransfers) == 0 {
		log.Fatalf("No successful backups available")
	}
	sort.Sort(TransfersList(completedValidTransfers))

	return completedValidTransfers[len(completedValidTransfers)-1]
}

func GetPublicUrl(t Transfer, herokuAppName string) PublicUrl {
	// TODO: dedup me into a client
	var herokuAuthToken = os.Getenv("HEROKU_AUTH_TOKEN")
	if herokuAuthToken == "" {
		log.Fatalf("HEROKU_AUTH_TOKEN must be set")
	}

	var publicUrlUrl = fmt.Sprintf(
		"https://postgres-api.heroku.com/client/v11/apps/%s/transfers/%d/actions/public-url",
		herokuAppName,
		t.Num,
	)

	client := &http.Client{}
	req, err := http.NewRequest("POST", publicUrlUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("", herokuAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	publicUrl := PublicUrl{}
	err = decoder.Decode(&publicUrl)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf(fmt.Sprintf("transfers call failed with %v", resp.StatusCode))
	}

	return publicUrl
}
