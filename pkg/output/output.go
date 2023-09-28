package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Output struct {
	FolderPath string
}

func NewOutput(folderPath string) *Output {
	return &Output{
		FolderPath: folderPath,
	}
}

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
		}
	}
}

// SaveJSON saves the OCR processing results as JSON to a file.
func (o *Output) SaveJSON(filename string, data interface{}) error {

	newDir := o.FolderPath
	err := os.MkdirAll(newDir, 0755)
	if err != nil {
		return err
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
			return err
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
		return err
	}

	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonFile.Write(formattedJson)
	jsonFile.Sync()

	return nil
}

// DisplayResults formats and displays the OCR processing results to the user.
func DisplayResults(text string) {
	// Display the extracted text
	fmt.Println("Extracted Text:")
	fmt.Println(text)

	// You can add more formatting and display options as needed
}
