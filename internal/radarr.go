package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var apiVersion string = "v3"

func GetMoviesData(url, b64ApiKey string) ([]byte, error) {
	apiUrl := url + "/api/" + apiVersion + "/movie"
	apiKey, err := Base64Decode(b64ApiKey)
	fmt.Println(apiKey)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("Failed to get movies library", err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to get movies library", err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Println("Failed to get movies library", res.Status)
		return nil, errors.New("Failed to get movies library")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to get movies library", err.Error())
		return nil, err
	}

	return data, nil
}
