package internal

import (
	"encoding/base64"
	"log"
	"strings"
)

func Base64Decode(encodedData string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		log.Println("Failed to decode base64 encoded data")
		return "", err
	}

	return strings.TrimSuffix(string(decodedData), "\n"), nil
}

func Base64Encode(decodedData string) string {
	encodedData := base64.StdEncoding.EncodeToString([]byte(decodedData))

	return encodedData
}
