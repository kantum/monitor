package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"monitor/internal/slack"
)

// Check if the website is up
func checkWebsite(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	return resp.StatusCode
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

	userIDs := os.Getenv("USER_IDS")
	if userIDs == "" {
		println("Missing SLACK_USER_IDS environment variable")
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
			for _, userID := range strings.Split(userIDs, ",") {
				msg += " <@" + userID + ">"
			}
			slack.SendSlackNotification(slackWebhook, msg)
			time.Sleep(15 * time.Minute)
		}

		// Wait for 5 seconds
		time.Sleep(5 * time.Second)
	}
}
