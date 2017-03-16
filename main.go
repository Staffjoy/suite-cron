// Service that queries the app repo cron repo on loop. app

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"papertrail"
	"time"
)

const (
	sleep                        = 40 // seconds
	consecutiveFailuresThreshold = 4
	papertrailPort               = 12345
	papertrailServer             = "logs1"
)

func main() {
	var env = os.Getenv("ENV")
	var url, apiKey string
	var logger *log.Logger

	switch env {
	case "stage":
		apiKey = os.Getenv("API_KEY")
		url = "https://stage.staffjoy.com"
		logger = log.New(
			&papertrail.Writer{
				Port:    papertrailPort,
				Network: papertrail.TCP,
				Server:  papertrailServer,
			}, "cron-stage ", log.LstdFlags)

	case "prod":
		apiKey = os.Getenv("API_KEY")
		url = "https://suite.staffjoy.com"
		logger = log.New(
			&papertrail.Writer{
				Port:    papertrailPort,
				Network: papertrail.TCP,
				Server:  papertrailServer,
			}, "cron-prod ", log.LstdFlags)

	default: // force all to dev
		env = "dev"
		url = "http://dev.staffjoy.com"
		apiKey = "staffjoydev"
		logger = log.New(os.Stdout, "cron-dev ", log.LstdFlags)
	}

	url += "/api/v2/internal/cron/"

	logger.Printf("INFO Initialized cron environment %s", env)

	// Boot a health check endpoint
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Cron running in env %s", env)
		})
		logger.Fatal(http.ListenAndServe(":80", nil))
	}()

	// configure request. stays same.
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Fatalf("ERROR - cannot initilize HTTP client - %v", err)
	}
	req.SetBasicAuth(apiKey, "")

	var consecutiveFailures int
	var success bool

	for {
		success = false
		start := time.Now()
		resp, err := client.Do(req)

		if err != nil {
			logger.Printf("INFO unable to query %s - %s", url, err)
		} else {
			duration := time.Since(start)

			switch resp.StatusCode {
			case http.StatusUnauthorized:
				logger.Printf("ERROR API key rejected on %s - %v", url, resp)
			case http.StatusTooManyRequests:
				logger.Printf("ERROR cron hitting rate limits on %s - %v", url, resp)
			case http.StatusInternalServerError:
				logger.Printf("ERROR cron api endpoint saying server error on %s - %v", url, resp)
			case http.StatusOK:
				success = true
				logger.Printf("INFO cron on %s took %.2f seconds", url, duration.Seconds())
			default:
				logger.Printf("INFO unhandled bad response from API %s - %v", url, resp)
			}
		}

		if success == true {
			consecutiveFailures = 0
		} else {
			consecutiveFailures++
			if consecutiveFailures >= consecutiveFailuresThreshold {
				logger.Printf("ERROR - unable to query %s after %d attempts ", url, consecutiveFailures)
			}
		}

		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
