package fileprocessing

import (
	"KlippaChallenge/pkg/api"
	"fmt"
	"os"
	"path/filepath"
)

func ProcessFile(filePath string, apiClient *api.Client, template string, pdfExtraction string, outputJSON string) (string, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("File %s not found", filePath)
	}

	// Call the OCR API to process the file
	text, err := apiClient.CallOCRService(filePath, template, pdfExtraction)
	if err != nil {
		return "", fmt.Errorf("Error processing file %s: %v", filePath, err)
	}

	// Save the JSON output if specified
	if outputJSON != "" {
		if err := saveJSONOutput(outputJSON, text); err != nil {
			return "", fmt.Errorf("Error saving JSON output: %v", err)
		}
	}

	return text, nil
}

// ProcessFolder processes all files in a folder using the provided API client.
func ProcessFolder(folderPath string, apiClient *api.Client, template string, pdfExtraction string, outputJSON string) {
	// Walk through the folder and process each file
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error walking through folder: %v\n", err)
			return err
		}
		if !info.IsDir() {
			// Process the file
			ProcessFile(path, apiClient, template, pdfExtraction, outputJSON)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error processing folder %s: %v\n", folderPath, err)
	}
}

// saveJSONOutput saves the extracted text to a JSON file.
func saveJSONOutput(outputJSON, text string) error {
	// Implement logic to save text to a JSON file here (e.g., using encoding/json)
	// If an error occurs during saving, return the error
	// Otherwise, return nil to indicate success
	return nil
}
