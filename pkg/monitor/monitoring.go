package monitor

import (
	"KlippaChallenge/pkg/api"
	"KlippaChallenge/pkg/fileprocessing"
	"KlippaChallenge/pkg/output"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Monitor struct {
	FolderPath    string
	APIClient     *api.APIClient
	Template      string
	PdfExtraction string
	OutputJSON    bool
}

func NewMonitor(folderPath string, apiClient *api.APIClient, template string, pdfExtraction string, outputJSON bool) *Monitor {
	return &Monitor{
		FolderPath:    folderPath,
		APIClient:     apiClient,
		Template:      template,
		PdfExtraction: pdfExtraction,
		OutputJSON:    outputJSON,
	}
}

func (m *Monitor) MonitorForNewFiles() error {
	processedFiles := make(map[string]bool)

	// Load processed files
	data, err := os.ReadFile("processed_files.log")
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			processedFiles[line] = true
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(m.FolderPath)
	if err != nil {
		return fmt.Errorf("error adding folder to watcher: %v", err)
	}

	fmt.Printf("Monitoring folder %s for new files...		(Ctrl+c to exit)\n", m.FolderPath)

	for event := range watcher.Events {
		if event.Op&fsnotify.Create == fsnotify.Create {
			if filepath.Ext(event.Name) == ".pdf" && !processedFiles[event.Name] {
				result, err := fileprocessing.ProcessFile(event.Name, m.APIClient, m.Template, m.PdfExtraction, m.OutputJSON)
				if err != nil {
					fmt.Printf("Error processing file %s: %v\n", event.Name, err)
				} else {
					fmt.Printf("Text extracted from %s\n", event.Name)
					fmt.Printf("Monitoring folder %s for new files...		(Ctrl+c to exit)\n", m.FolderPath)
					outputClient := output.NewOutput("json_outputs")
					if m.OutputJSON {
						outputClient.SaveJSON(event.Name, result)
					} else {
						output.DisplayResults(result)
					}

					// Mark as processed
					processedFiles[event.Name] = true
					func() {
						f, err := os.OpenFile("processed_files.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
						if err != nil {
							log.Fatalf("Failed to open log file: %v", err)
						}
						defer f.Close()
						if _, err = f.WriteString(event.Name + "\n"); err != nil {
							log.Fatalf("Failed to write to log file: %v", err)
						}
					}()
				}
			}
		}
	}

	return nil
}
