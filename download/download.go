package download

import (
	"io"
	"log"
	"net/http"
	"os"
)

// DownloadUrl takes a url string, downloads it, and copies it to a local file
// on disk.
func DownloadUrl(url string) (*os.File, error) {
	out, err := os.Create("backup")
	defer out.Close()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Copied %d bytes", n)

	return out, nil
}
