// Package output provides functions for outputting results.
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Output represents an output destination.
type Output struct {
	FolderPath string
}

// NewOutput creates a new Output with the provided folder path.
func NewOutput(folderPath string) *Output {
	return &Output{
		FolderPath: folderPath,
	}
}

// deleteEmpty recursively deletes all empty values from a map.
func (o *Output) deleteEmpty(mapData map[string]interface{}) {
	for key, value := range mapData {
		switch v := value.(type) {
		case string:
			if v == "" {
				delete(mapData, key)
			}
		case nil:
			delete(mapData, key)
		case map[string]interface{}:
			o.deleteEmpty(v)
			if len(v) == 0 { // If the nested map is empty after deletion, remove it from the map
				delete(mapData, key)
			}
		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					o.deleteEmpty(itemMap)
					if len(itemMap) == 0 { // If the nested map is empty after deletion, remove it from the array
						v = append(v[:i], v[i+1:]...)
					}
				}
			}
			if len(v) == 0 {
				delete(mapData, key)
			} else {
				mapData[key] = v // Update the array in the original map
			}
		default:
			// For all other data types, do nothing.
		}
	}
}

// SaveJSON saves the OCR processing results as JSON to a file.
func (o *Output) SaveJSON(filename string, data interface{}) error {

	newDir := o.FolderPath
	err := os.MkdirAll(newDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	baseName := filepath.Base(filename)

	filename = filepath.Join(newDir, baseName)
	ext := filepath.Ext(filename)
	jsonFileName := filename[0:len(filename)-len(ext)] + ".json"

	var jsonData map[string]interface{}

	if dataMap, ok := data.(map[string]interface{}); ok {
		jsonData = dataMap
	} else if dataStr, ok := data.(string); ok {
		err := json.Unmarshal([]byte(dataStr), &jsonData)
		if err != nil {
			return fmt.Errorf("error unmarshalling JSON: %w", err)
		}
	} else {
		return fmt.Errorf("invalid data type for SaveJSON")
	}

	dataMap, ok := jsonData["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error: 'data' field is not a JSON object")
	}

	o.deleteEmpty(dataMap)

	formattedJson, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer jsonFile.Close()

	jsonFile.Write(formattedJson)
	jsonFile.Sync()

	return nil
}

// DisplayResults formats and displays the OCR processing results to the user.
func DisplayResults(text string) {
	fmt.Println("Extracted Text:")
	fmt.Println(text)
}
