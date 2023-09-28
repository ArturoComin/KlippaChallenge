package main

import (
	"KlippaChallenge/pkg/api"
	"KlippaChallenge/pkg/monitor"
	"KlippaChallenge/pkg/output"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

type Config struct {
	APIKey         string
	Template       string
	ExtractionType string
	OutputJSON     bool
	DirPath        string
	Monitoring     bool
}

func main() {

	configpath := flag.String("config", "config.json", "Path to file containing the configuration")
	config, err := LoadConfig(configpath)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	flag.StringVar(&config.APIKey, "api-key", "config.APIKey", "API key for the OCR API")
	flag.StringVar(&config.Template, "template", "config.Template", "Template for the OCR API")
	flag.StringVar(&config.ExtractionType, "pdf-extraction", "config.ExtractionType", "PDF text extraction type: fast or full")
	flag.BoolVar(&config.OutputJSON, "output-json", false, "Set to true to get JSON output")
	flag.StringVar(&config.DirPath, "dir-path", "", "Path to directory containing PDF files")
	flag.BoolVar(&config.Monitoring, "monitoring", false, "Set to true to monitor the directory")

	flag.Parse()

	if config.APIKey == "" || config.Template == "" {
		fmt.Println("Please provide all the required options: api-key, template")
		os.Exit(1)
	}

	apiClient := api.NewAPIClient(config.APIKey)

	if config.DirPath != "" {
		files, err := os.ReadDir(config.DirPath)
		if err != nil {
			panic(err)
		}
		var wg sync.WaitGroup
		for _, file := range files {
			wg.Add(1)
			go func(file os.DirEntry) {
				defer wg.Done()
				filePath := path.Join(config.DirPath, file.Name())
				result, err := apiClient.CallAPI(config.Template, config.ExtractionType, filePath)
				if err != nil {
					panic(err)
				}
				outputClient := output.NewOutput("json_outputs")
				if config.OutputJSON {
					outputClient.SaveJSON(filePath, result)
				} else {
					output.DisplayResults(result)
				}
			}(file)
		}
		wg.Wait()
	} else {

		if flag.NArg() == 0 {
			fmt.Println("Usage: my-cli-tool -api-key <API_KEY> -template <TEMPLATE> -pdf-extraction <MODE> -output-json <JSON_FILE> <FILE>")
			os.Exit(1)
		}

		filePath := flag.Arg(0)

		result, err := apiClient.CallAPI(config.Template, config.ExtractionType, filePath)
		if err != nil {
			panic(err)
		}

		outputClient := output.NewOutput("json_outputs")
		if config.OutputJSON {
			outputClient.SaveJSON(filePath, result)
		} else {
			output.DisplayResults(result)
		}
	}

	monitor := monitor.NewMonitor(config.DirPath, apiClient, config.Template, config.ExtractionType, config.OutputJSON)
	if config.Monitoring {
		go monitor.MonitorForNewFiles()
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
	}

}

func LoadConfig(configpath *string) (*Config, error) {

	config := &Config{}
	file, err := os.Open(*configpath)
	if os.IsNotExist(err) {
		return config, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
