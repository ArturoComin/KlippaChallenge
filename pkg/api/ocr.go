package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type APIClient struct {
	Client *http.Client
	APIKey string
}

func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		Client: &http.Client{},
		APIKey: apiKey,
	}
}

func (c *APIClient) CallAPI(template string, extractionType string, filePath string) (string, error) {
	request, err := http.NewRequest("POST", "https://custom-ocr.klippa.com/api/v1/parseDocument", nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("X-Auth-Key", c.APIKey)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("document", filePath)
	if err != nil {
		return "", err
	}
	pdfBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	_, err = part.Write(pdfBytes)
	if err != nil {
		return "", err
	}

	err = writer.WriteField("pdf_text_extraction", extractionType)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Body = io.NopCloser(body)

	response, err := c.Client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
