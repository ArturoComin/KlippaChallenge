// Package fileprocessing provides functions for processing files.
package fileprocessing

import (
	"KlippaChallenge/pkg/api"
	"KlippaChallenge/pkg/output"
	"fmt"
	"os"
	"path"
	"sync"
)

// ProcessFile processes a single file.
func ProcessFile(filename string, apiClient *api.APIClient, template string, pdfExtraction string, outputJSON bool) error {
	// Print a message to indicate that the file is being processed.
	fmt.Printf("Processing file: %s\n", filename)

	// Call the API with the file.
	result, err := apiClient.CallAPI(template, pdfExtraction, filename)
	if err != nil {
		return fmt.Errorf("error calling API: %w", err)
	}

	// Create a new output client and save or display the results.
	outputClient := output.NewOutput("json_outputs")
	if outputJSON {
		outputClient.SaveJSON(filename, result)
	} else {
		output.DisplayResults(result)
	}

	return nil
}

// ProcessFolder processes all files in a folder concurrently.
func ProcessFolder(dirPath string, apiClient *api.APIClient, template string, extractionType string, outputJSON bool) error {
	// Read the directory.
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			filePath := path.Join(dirPath, file.Name())

			result, err := apiClient.CallAPI(template, extractionType, filePath)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", filePath, err)
				return
			}

			outputClient := output.NewOutput("json_outputs")
			if outputJSON {
				outputClient.SaveJSON(filePath, result)
			} else {
				output.DisplayResults(result)
			}
		}(file)
	}
	wg.Wait()

	return nil
}
