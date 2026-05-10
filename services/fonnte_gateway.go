package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// FonnteWAGateway sends WhatsApp messages via the Fonnte.com API.
type FonnteWAGateway struct {
	apiURL string
	token  string
	client *http.Client
}

// NewFonnteWAGateway creates a new Fonnte gateway with the given API URL and token.
func NewFonnteWAGateway(apiURL, token string) *FonnteWAGateway {
	return &FonnteWAGateway{
		apiURL: apiURL,
		token:  token,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// SendMessage sends a WhatsApp text message to the target phone number via Fonnte.
func (f *FonnteWAGateway) SendMessage(target, message string) error {
	payload := map[string]string{
		"target":      target,
		"message":     message,
		"countryCode": "62", // Indonesia
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fonnte marshal: %w", err)
	}

	req, err := http.NewRequest("POST", f.apiURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("fonnte request: %w", err)
	}

	req.Header.Set("Authorization", f.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("fonnte send: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Fonnte returns 200 for both success and some errors, parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("fonnte response parse error: %s", string(respBody))
	}

	// Check if the status indicates success
	if status, ok := result["status"].(bool); ok && !status {
		reason, _ := result["reason"].(string)
		return fmt.Errorf("fonnte error: %s", reason)
	}

	log.Printf("📱 [FONNTE] Message sent to %s — response: %s", target, string(respBody))
	return nil
}
