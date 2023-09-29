// Package main provides the entry point for the application.
package main

import (
	"KlippaChallenge/pkg/api"
	"KlippaChallenge/pkg/fileprocessing"
	"KlippaChallenge/pkg/monitor"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Config represents the configuration for the application.
type Config struct {
	APIKey         string
	Template       string
	ExtractionType string
	OutputJSON     bool
	DirPath        string
	Monitoring     bool
}

func main() {

	// Load configuration from config file provided by user
	configpath := flag.String("config", "config.json", "Path to file containing the configuration")
	config, err := LoadConfig(configpath)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	// Parse flags to config
	flag.StringVar(&config.APIKey, "api-key", config.APIKey, "API key for the OCR API")
	flag.StringVar(&config.Template, "template", config.Template, "Template for the OCR API")
	flag.StringVar(&config.ExtractionType, "pdf-extraction", config.ExtractionType, "PDF text extraction type: fast or full")
	flag.BoolVar(&config.OutputJSON, "output-json", config.OutputJSON, "Set to true to get JSON output")
	flag.StringVar(&config.DirPath, "dir-path", config.DirPath, "Path to directory containing PDF files")
	flag.BoolVar(&config.Monitoring, "monitoring", config.Monitoring, "Set to true to monitor the directory")

	flag.Parse()

	apiClient := api.NewAPIClient(config.APIKey)

	// Check if the config is valid
	if config.APIKey == "" || config.Template == "" {
		fmt.Println("Please provide all the required options: -api-key, -template")
		os.Exit(1)
	}

	// Process PDF files if in directory
	if config.DirPath != "" {
		err = fileprocessing.ProcessFolder(config.DirPath, apiClient, config.Template, config.ExtractionType, config.OutputJSON)
		if err != nil {
			log.Fatalf("Failed to process folder: %s", err)
		}
	} else {

		if flag.NArg() == 0 {
			fmt.Println("Please provide a file or directory path using the -dir-path flag")
			os.Exit(1)
		}

		filePath := flag.Arg(0)
		err = fileprocessing.ProcessFile(filePath, apiClient, config.Template, config.ExtractionType, config.OutputJSON)
		if err != nil {
			log.Fatalf("Failed to process file: %s", err)
		}
	}

	fileMonitor := monitor.NewMonitor(config.DirPath, apiClient, config.Template, config.ExtractionType, config.OutputJSON)
	if config.Monitoring {
		go func() {
			err = fileMonitor.MonitorForNewFiles()
			if err != nil {
				log.Fatalf("Failed to monitor for new files: %s", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
	}
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfig(configpath *string) (*Config, error) {

	config := &Config{}
	file, err := os.Open(*configpath)
	if os.IsNotExist(err) {
		// If the file does not exist, return an empty configuration.
		return config, nil
	} else if err != nil {
		// If there was an error opening the file (other than it not existing), return the error.
		return nil, fmt.Errorf("error opening configuration file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		// If there was an error decoding the JSON file into a Config struct,
		// return an error describing what went wrong.
		return nil, fmt.Errorf("error decoding configuration file: %w", err)
	}

	return config, nil
}
