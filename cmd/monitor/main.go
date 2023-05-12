package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// Create a struct to read the JSON body
type slackBody struct {
	Text string `json:"text"`
}

// Check if the website is up
func checkWebsite(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	return resp.StatusCode
}

// Send a POST request to a Slack webhook
func sendSlackNotification(webhook string, message string) {
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

func main() {
	slackWebhook := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhook == "" {
		println("Missing SLACK_WEBHOOK_URL environment variable")
		return
	}
	websiteURL := os.Getenv("WEBSITE_URL")
	if websiteURL == "" {
		println("Missing WEBSITE_URL environment variable")
		return
	}

	for {
		// Check if the website is up
		httpStatus := checkWebsite(websiteURL)
		switch httpStatus {
		case 0:
			println("No connection")
		case 200:
			println("Website is up")
		case 404:
			println("Website returned 404 - not found")
		case 500:
			println("Website returned 500 - internal server error")
		default:
			println("Website returned ", httpStatus)
		}

		if httpStatus != 200 && httpStatus != 0 {
			println("Website is down")
			msg := "SmartCodec is down, next check in 15 minutes."
			sendSlackNotification(slackWebhook, msg)
			time.Sleep(15 * time.Minute)
		}

		// Wait for 5 seconds
		time.Sleep(5 * time.Second)
	}
}
