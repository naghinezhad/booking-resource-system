package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultTotalRequests = 1000
	defaultAPIURL        = "http://localhost:8080/reserve"
	defaultResourceID    = "665e3c4a0f7b8b3c8c6a1234"
)

type reserveRequest struct {
	ResourceID string `json:"resource_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

func main() {
	totalRequests := int64(getIntEnv("LOAD_TEST_TOTAL_REQUESTS", defaultTotalRequests))
	apiURL := getStringEnv("LOAD_TEST_API_URL", defaultAPIURL)
	resourceID := getStringEnv("LOAD_TEST_RESOURCE_ID", defaultResourceID)

	var (
		success     int64 // 201 Created
		conflicts   int64 // 409 Conflict (double booking prevented)
		rateLimited int64 // 429 Too Many Requests (queue full / limiter)
		failed      int64 // 500 or network errors
	)

	start := time.Now()
	var wg sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(1 * time.Hour)

	payload := reserveRequest{
		ResourceID: resourceID,
		StartTime:  startTime.Format(time.RFC3339),
		EndTime:    endTime.Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("failed to marshal request body: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting load test: %d requests\n", totalRequests)
	fmt.Printf("Target API: %s\n", apiURL)
	fmt.Printf("Resource ID: %s\n", resourceID)

	for i := int64(0); i < totalRequests; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(body))
			if err != nil {
				atomic.AddInt64(&failed, 1)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				atomic.AddInt64(&failed, 1)
				return
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusCreated:
				atomic.AddInt64(&success, 1)
			case http.StatusConflict:
				atomic.AddInt64(&conflicts, 1)
			case http.StatusTooManyRequests:
				atomic.AddInt64(&rateLimited, 1)
			default:
				atomic.AddInt64(&failed, 1)
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("========================================")
	fmt.Println("           LOAD TEST RESULTS")
	fmt.Println("========================================")
	fmt.Printf("Total Requests:      %d\n", totalRequests)
	fmt.Printf("Duration:            %v\n", duration)
	fmt.Printf("Requests/sec:        %.2f\n", float64(totalRequests)/duration.Seconds())
	fmt.Println("----------------------------------------")
	fmt.Printf("Success (201):       %d (expected 1)\n", success)
	fmt.Printf("Conflicts (409):     %d (double booking prevented)\n", conflicts)
	fmt.Printf("Rate Limited (429):  %d (queue/limiter dropped)\n", rateLimited)
	fmt.Printf("Failed (other):      %d\n", failed)
	fmt.Println("========================================")
}

func getStringEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getIntEnv(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
