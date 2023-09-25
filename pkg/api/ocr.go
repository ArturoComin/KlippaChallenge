package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents the Klippa OCR API client.
type Client struct {
	apiKey  string
	baseURL string
}

// NewClient creates a new API client with the provided API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://custom-ocr.klippa.com/api/v1", // Adjust the base URL accordingly
	}
}

// OCRResponse represents the response from the OCR API.
type OCRResponse struct {
	// Define the fields you expect in the API response
	Text string `json:"text"`
	// Add more fields as needed
}

// ProcessFile sends a file for OCR processing and returns the extracted text.
func (c *Client) CallOCRService(filePath string, template string, pdfExtraction string) (string, error) {
	// Create an HTTP client
	client := &http.Client{}

	// Prepare the API request payload
	requestData := map[string]interface{}{
		"template":         template,
		"pdf_text_extract": pdfExtraction,
		// Add other request parameters here
	}

	// Convert the request payload to JSON
	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Create a request to the Klippa OCR API
	req, err := http.NewRequest("POST", c.baseURL+"/ocr", bytes.NewBuffer(requestDataJSON))
	if err != nil {
		return "", err
	}

	// Set the API key header
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the API response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read and parse the API response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal the response JSON into the OCRResponse struct
	var ocrResponse OCRResponse
	if err := json.Unmarshal(responseBody, &ocrResponse); err != nil {
		return "", err
	}

	// Extracted text from the API response
	return ocrResponse.Text, nil
}
