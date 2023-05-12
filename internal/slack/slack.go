package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Create a struct to read the JSON body
type slackBody struct {
	Text string `json:"text"`
}

// Send a POST request to a Slack webhook
func SendSlackNotification(webhook string, message string) {
	// Create the request body struct
	slackBody := slackBody{Text: message}

	// Marshal the struct into JSON
	slackBodyBytes, _ := json.Marshal(slackBody)

	// Send a POST request with the JSON as the body
	req, err := http.NewRequest(
		http.MethodPost,
		webhook,
		bytes.NewBuffer(slackBodyBytes),
	)
	if err != nil {
		return
	}

	// Set the request header Content-Type for JSON
	req.Header.Add("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
