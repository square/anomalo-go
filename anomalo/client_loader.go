package anomalo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const DefaultAnomaloSecretsFile = "anomalo_secrets.json"

// CreateClientFromFile Creates a client based on credentials in a local file.
func CreateClientFromFile(filePath string) (*Client, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	contents, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var client Client
	err = json.Unmarshal(contents, &client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// CreateClientFromEnv Creates a client based on credentials in environment
// variables.
func CreateClientFromEnv() (*Client, error) {
	token := os.Getenv("ANOMALO_API_SECRET_TOKEN")
	host := os.Getenv("ANOMALO_INSTANCE_HOST")
	if token == "" || host == "" {
		return nil, fmt.Errorf(
			"at least one anomalo API env variable is not set. Got host '%s' and token length %d",
			host, len(token), // Don't log the token
		)
	}
	client := Client{
		Token: token,
		Host:  host,
	}
	return &client, nil
}

// CreateClient Attempts to instantiate a Client. First checks for Anomalo
// credentials in anomalo_secrets.json, then checks in environment variables.
func CreateClient() (*Client, error) {
	var client *Client
	client, err := CreateClientFromFile(DefaultAnomaloSecretsFile)
	if err != nil {
		log.Printf("Did not find local anomalo credentials in %s. Checking environment variables.\n", DefaultAnomaloSecretsFile)
		client, err = CreateClientFromEnv()
		if err != nil {
			log.Printf(err.Error())
			return nil, fmt.Errorf("could not find anomalo credentials")
		}
	}

	return client, nil
}
