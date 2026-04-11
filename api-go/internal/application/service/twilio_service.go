package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// TwilioClient provides integration with Twilio SMS API
type TwilioClient struct {
	AccountSID  string
	AuthToken   string
	PhoneNumber string
	BaseURL     string
	httpClient  *http.Client
}

// TwilioMessageResponse represents the response from Twilio API
type TwilioMessageResponse struct {
	SID         string `json:"sid"`
	AccountSID  string `json:"account_sid"`
	To          string `json:"to"`
	From        string `json:"from"`
	Status      string `json:"status"`
	Body        string `json:"body"`
	DateCreated string `json:"date_created"`
	ErrorCode   *int   `json:"error_code"`
	ErrorMessage *string `json:"error_message"`
	Price       string `json:"price"`
}

// TwilioErrorResponse represents an error response from Twilio
type TwilioErrorResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	MoreInfo    string `json:"more_info"`
	Status      int    `json:"status"`
}

// NewTwilioClient creates a new Twilio client from environment variables
func NewTwilioClient() *TwilioClient {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	phoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	if accountSID == "" || authToken == "" {
		return nil // Twilio not configured
	}

	return &TwilioClient{
		AccountSID:  accountSID,
		AuthToken:   authToken,
		PhoneNumber: phoneNumber,
		BaseURL:     "https://api.twilio.com/2010-04-01",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends an SMS message via Twilio
func (c *TwilioClient) SendMessage(to, message string) (messageSID string, cost float64, err error) {
	if c == nil {
		return "", 0, fmt.Errorf("twilio client not configured")
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("To", formatPhoneNumber(to))
	formData.Set("From", c.PhoneNumber)
	formData.Set("Body", message)

	// Build URL
	u := fmt.Sprintf("%s/Accounts/%s/Messages.json", c.BaseURL, c.AccountSID)

	// Create request
	req, err := http.NewRequest("POST", u, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.AccountSID, c.AuthToken)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var twilioErr TwilioErrorResponse
		if err := json.Unmarshal(body, &twilioErr); err == nil {
			return "", 0, fmt.Errorf("twilio error %d: %s", twilioErr.Code, twilioErr.Message)
		}
		return "", 0, fmt.Errorf("twilio HTTP error %d: %s", resp.StatusCode, string(body))
	}

	// Parse successful response
	var msgResp TwilioMessageResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w", err)
	}

	// Parse cost
	cost = 0.0
	if msgResp.Price != "" {
		fmt.Sscanf(msgResp.Price, "%f", &cost)
		if strings.HasPrefix(msgResp.Price, "-") {
			cost = -cost // Twilio returns negative for charges
		}
	}

	return msgResp.SID, cost, nil
}

// GetMessageStatus retrieves the status of a sent message
func (c *TwilioClient) GetMessageStatus(messageSID string) (*TwilioMessageResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("twilio client not configured")
	}

	u := fmt.Sprintf("%s/Accounts/%s/Messages/%s.json", c.BaseURL, c.AccountSID, messageSID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.AccountSID, c.AuthToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get message status: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("twilio HTTP error %d: %s", resp.StatusCode, string(body))
	}

	var msgResp TwilioMessageResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &msgResp, nil
}

// formatPhoneNumber ensures the phone number is in E.164 format
func formatPhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)

	// If already starts with +, return as-is
	if strings.HasPrefix(phone, "+") {
		return phone
	}

	// Remove any non-digit characters
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// If starts with 1 and 11 digits, it's US number
	if len(digits) == 11 && strings.HasPrefix(digits, "1") {
		return "+" + digits
	}

	// If 10 digits, assume US number
	if len(digits) == 10 {
		return "+1" + digits
	}

	// Return with + prefix
	return "+" + digits
}
