// Package api provides a client for making API requests.
package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// APIClient is a client for making API requests.
type APIClient struct {
	Client *http.Client
	APIKey string
}

// NewAPIClient creates a new APIClient with the provided API key.
func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		Client: &http.Client{},
		APIKey: apiKey,
	}
}

// CallAPI makes a POST request to the Klippa OCR API with the provided parameters.
func (c *APIClient) CallAPI(template string, extractionType string, filePath string) (string, error) {
	// Create a new POST request.
	request, err := http.NewRequest("POST", "https://custom-ocr.klippa.com/api/v1/parseDocument", nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Add("X-Auth-Key", c.APIKey)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the file to the request body.
	part, err := writer.CreateFormFile("document", filePath)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}
	pdfBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	_, err = part.Write(pdfBytes)
	if err != nil {
		return "", fmt.Errorf("error writing file to form: %w", err)
	}

	err = writer.WriteField("pdf_text_extraction", extractionType)
	if err != nil {
		return "", fmt.Errorf("error writing field to form: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Body = io.NopCloser(body)

	response, err := c.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	return string(bodyBytes), nil
}
