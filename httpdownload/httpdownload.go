package httpdownload

import (
	"io"
	"log"
	"net/http"
	"os"
)

func HTTPDownload(url string, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return err
	}
	defer outFile.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error downloading file: %s", err)
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(outFile, resp.Body)
	return err
}
