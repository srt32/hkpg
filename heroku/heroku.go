package heroku

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
)

// Client allows for http calls to Heroku
type Client struct {
	HTTP *http.Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) NewRequest(method, path string) (*http.Request, error) {
	var herokuAuthToken = os.Getenv("HEROKU_AUTH_TOKEN")
	if herokuAuthToken == "" {
		log.Fatalf("HEROKU_AUTH_TOKEN must be set")
	}

	const apiURL = "https://postgres-api.heroku.com"

	req, err := http.NewRequest(method, apiURL+path, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("", herokuAuthToken)
	return req, nil
}

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
	var transfersPath = fmt.Sprintf(
		"/client/v11/apps/%s/transfers",
		appName,
	)

	client := Client{}
	req, err := client.NewRequest("GET", transfersPath)
	if err != nil {
		log.Fatal(err)
	}

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
	var publicUrlPath = fmt.Sprintf(
		"/client/v11/apps/%s/transfers/%d/actions/public-url",
		herokuAppName,
		t.Num,
	)

	client := Client{}
	req, err := client.NewRequest("POST", publicUrlPath)
	if err != nil {
		log.Fatal(err)
	}

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
		log.Fatalf(fmt.Sprintf("public url call failed with %v", resp.StatusCode))
	}

	return publicUrl
}
