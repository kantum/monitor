package main

import (
	"fmt"
	"monitor/internal/slack"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Check if the website is up
func checkWebsite(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	return resp.StatusCode
}

// printVCSInfo prints the selected VCS info
func printVCSInfo() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" || setting.Key == "vcs.time" {
				fmt.Printf("%+v\t%+v\n", setting.Key, setting.Value)
			}
		}
	}
}

func main() {
	// Print the VCS info
	printVCSInfo()

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

		// If the website is down, send a notification and wait for 15 minutes
		if httpStatus != 200 && httpStatus != 0 {
			println("Website is down")
			msg := "SmartCodec is down, next check in 15 minutes."
			for _, userID := range strings.Split(userIDs, ",") {
				msg += " <@" + userID + ">"
			}
			slack.PostWebhook(slackWebhook, msg)
			time.Sleep(15 * time.Minute)
		}

		// Wait for 5 seconds
		time.Sleep(5 * time.Second)
	}
}
