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

// Do passes HTTP requests to the underlying client and executes them
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

// NewRequest creates an HTTP request that should be passed to *Client.Do
func (c *Client) NewRequest(method, path string) (*http.Request, error) {
	var herokuAuthToken = os.Getenv("HKPG_HEROKU_AUTH_TOKEN")
	if herokuAuthToken == "" {
		log.Fatalf("HKPG_HEROKU_AUTH_TOKEN must be set")
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

// GetTransfers takes a Heroku app name and returns the most recent successful
// backup
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
	if resp.StatusCode != 200 {
		log.Fatalf(fmt.Sprintf("transfers call failed with %v", resp.StatusCode))
	}

	var transfers = make([]Transfer, 0)
	err = decoder.Decode(&transfers)
	if err != nil {
		log.Fatal(err)
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

// GetPublicUrl takes a Transfer and app name and returns the URL of a publicly
// accessible signed URL for the most recent backup.
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
