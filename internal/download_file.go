package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(url, dest string) error {

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Failed to download file to", dest, "Error:", resp.Status)
		return errors.New("Failed to download file")
	}

	file, err := os.Create(dest)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("File downloaded.")

	return nil
}
