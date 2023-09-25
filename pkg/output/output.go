package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// DisplayResults formats and displays the OCR processing results to the user.
func DisplayResults(text string) {
	// Display the extracted text
	fmt.Println("Extracted Text:")
	fmt.Println(text)

	// You can add more formatting and display options as needed
}

// SaveJSON saves the OCR processing results as JSON to a file.
func SaveJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = writeFile(filename, jsonData)
	return err
}

// Helper function to write data to a file
func writeFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	fmt.Printf("Results saved to %s\n", filename)
	return nil
}
