package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time" // Added for polling delay

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require" // Added for fatal assertions
)

func TestServer(t *testing.T) {
	go startServer()

	var resp *http.Response
	var err error
	const maxRetries = 20
	const retryDelay = 250 * time.Millisecond

	serverReady := false
	for i := 0; i < maxRetries; i++ {
		healthResp, httpErr := http.Get("http://localhost:8080/health")
		if httpErr == nil {
			if healthResp.Body != nil {
				healthResp.Body.Close()
			}
			if healthResp.StatusCode == http.StatusOK {
				serverReady = true
				break
			}
		}
		time.Sleep(retryDelay)
	}

	require.True(t, serverReady, "Server did not become ready in time. Last error: %v", err)

	resp, err = http.Get("http://localhost:8080/")
	require.NoError(t, err, "Error getting base path /")
	require.NotNil(t, resp, "Response from base path / was nil")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Base path / did not return HTTP 200 OK")

	respExecution, err := http.Get("http://localhost:8080/execute?cep=04475030")
	require.NoError(t, err, "Error getting /execute with cep=0447503")
	require.NotNil(t, respExecution, "Response from /execute was nil")
	defer respExecution.Body.Close()
	assert.Equal(t, http.StatusOK, respExecution.StatusCode, "/execute did not return HTTP 200 OK")
	if respExecution.StatusCode != http.StatusOK {
		t.Fatalf("Expected HTTP 200 OK, got %d", respExecution.StatusCode)
	}

	dataResponse := SuccessResponse{}

	err = json.NewDecoder(respExecution.Body).Decode(&dataResponse)
	if err != nil {
		t.Fatalf("Error decoding JSON response from /execute: %v", err)
	}

	t.Logf("Response from /execute: %+v", dataResponse)
}
