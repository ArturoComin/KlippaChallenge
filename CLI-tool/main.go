package main

import (
	"KlippaChallenge/pkg/api"
	"KlippaChallenge/pkg/config"
	"KlippaChallenge/pkg/fileprocessing"
	"KlippaChallenge/pkg/output"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define command-line flags and options
	apiKey := flag.String("api-key", "", "Klippa OCR API key")
	template := flag.String("template", "", "OCR template")
	pdfExtraction := flag.String("pdf-extraction", "fast", "PDF text extraction mode (fast or full)")
	outputJSON := flag.String("output-json", "", "JSON output file name")
	folderPath := flag.String("folder", "", "Path to a folder for batch processing")

	flag.Parse()

	// Load configuration (optional)
	config.LoadConfig() // You can implement this function to load global configuration from a file

	// Check if the API key is provided
	if *apiKey == "" {
		fmt.Println("API key is required. Use the -api-key flag.")
		os.Exit(1)
	}

	// Initialize API client
	apiClient := api.NewClient(*apiKey)

	// Process files or folder
	if *folderPath != "" {
		// Batch processing of a folder
		fileprocessing.ProcessFolder(*folderPath, apiClient, *template, *pdfExtraction, *outputJSON)
	} else {
		// Single file processing
		if flag.NArg() == 0 {
			fmt.Println("Usage: my-cli-tool -api-key <API_KEY> -template <TEMPLATE> -pdf-extraction <MODE> -output-json <JSON_FILE> <FILE>")
			os.Exit(1)
		}
		inputFile := flag.Arg(0)

		// Call fileprocessing.ProcessFile
		text, err := fileprocessing.ProcessFile(inputFile, apiClient, *template, *pdfExtraction, *outputJSON)
		if err != nil {
			fmt.Printf("Error processing file: %v\n", err)
			os.Exit(1)
		}
		// Display the extracted text
		output.DisplayResults(text)

	}

}
