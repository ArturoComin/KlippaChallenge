// Package monitor provides a Monitor for watching a folder for new files.
package monitor

import (
	"KlippaChallenge/pkg/api"
	"KlippaChallenge/pkg/fileprocessing"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// Monitor watches a folder for new files and processes them.
type Monitor struct {
	FolderPath    string
	APIClient     *api.APIClient
	Template      string
	PdfExtraction string
	OutputJSON    bool
}

// NewMonitor creates a new Monitor with the provided parameters.
func NewMonitor(folderPath string, apiClient *api.APIClient, template string, pdfExtraction string, outputJSON bool) *Monitor {
	return &Monitor{
		FolderPath:    folderPath,
		APIClient:     apiClient,
		Template:      template,
		PdfExtraction: pdfExtraction,
		OutputJSON:    outputJSON,
	}
}

// MonitorForNewFiles watches the Monitor's folder for new files and processes them.
func (m *Monitor) MonitorForNewFiles() error {
	processedFiles := make(map[string]bool)

	// Load processed files
	data, err := os.ReadFile("processed_files.log")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading processed files log: %w", err)
	}
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if len(line) > 0 {
				processedFiles[line] = true
			}
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer watcher.Close()

	err = watcher.Add(m.FolderPath)
	if err != nil {
		return fmt.Errorf("error adding folder to watcher: %w", err)
	}

	fmt.Printf("Monitoring folder %s for new files... (Ctrl+c to exit)\n", m.FolderPath)

	var wg sync.WaitGroup
	for event := range watcher.Events {
		if event.Op&fsnotify.Create == fsnotify.Create {
			if filepath.Ext(event.Name) == ".pdf" && !processedFiles[event.Name] {
				wg.Add(1)
				go func(eventName string) {
					defer wg.Done()
					err = fileprocessing.ProcessFile(eventName, m.APIClient, m.Template, m.PdfExtraction, m.OutputJSON)
					if err != nil {
						fmt.Printf("Error processing file %s: %v\n", eventName, err)
					}
				}(event.Name)
			}
		}
	}

	wg.Wait()

	return nil
}
