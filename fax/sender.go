package fax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	fax_response "github.com/kvizdos/notifyre-client/fax/responses"
)

func Send(payload Payload, apiKey string) (*fax_response.SuccessResponse, error) {
	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", "https://api.notifyre.com/fax/send", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set headers
	req.Header = http.Header{
		"x-api-token":  {apiKey},
		"Content-Type": {"application/json"},
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != 200 {
		var apiResponse fax_response.FailureResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to send request unmarshal failure: %s", err.Error())
		}
		return nil, apiResponse
	}

	var apiResponse fax_response.SuccessResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to send request unmarshal failure: %s", err.Error())
	}

	return &apiResponse, nil
}
