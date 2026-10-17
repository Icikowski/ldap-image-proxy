// Package main provides the container healthcheck probe.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("LIP_SERVER_PORT")
	if port == "" {
		port = "3891"
	}

	if err := executeChecks(port); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func executeChecks(port string) error {
	var (
		client  = &http.Client{Timeout: 3 * time.Second}
		baseURL = "http://127.0.0.1:" + port
	)

	return probe(client, baseURL)
}

func probe(client *http.Client, url string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		dataStr := string(data)
		if dataStr == "" {
			dataStr = "(no response body)"
		}

		//nolint:err113
		return fmt.Errorf("%s returned %s: %s", url, resp.Status, dataStr)
	}

	return nil
}
