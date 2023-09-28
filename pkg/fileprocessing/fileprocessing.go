package fileprocessing

import (
	"KlippaChallenge/pkg/api"
	"fmt"
)

func ProcessFile(filename string, apiClient *api.APIClient, template string, pdfExtraction string, outputJSON bool) (string, error) {
	// Process the file here...
	fmt.Printf("Processing file: %s\n", filename)

	// For example, call the API with the file
	result, err := apiClient.CallAPI(template, pdfExtraction, filename)
	if err != nil {
		return "", fmt.Errorf("error calling API: %w", err)
	}

	return result, nil
}
