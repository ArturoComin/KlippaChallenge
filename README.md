# Klippa OCR Challenge

This is a Go application that uses the Klippa OCR API to process PDF files. It can process a single file, all files in a directory, or monitor a directory for new files and process them as they are added.

## Features

- Process a single PDF file
- Process all PDF files in a directory
- Monitor a directory for new files and process them as they are added
- Output results as JSON
- Runs in Docker

## Usage

You can provide options to the application using command line flags or a configuration file. The available options are:

- `-api-key`: The API key for the OCR API.
- `-template`: The template for the OCR API.
- `-pdf-extraction`: The PDF text extraction type: `fast` or `full`.
- `-output-json`: Set to true to get JSON output.
- `-dir-path`: The path to the directory containing PDF files.
- `-monitoring`: Set to true to monitor the directory.

For example, you can run the application with command line flags like this:

>go run main.go -api-key your-api-key -template your-template -pdf-extraction fast -output-json true -dir-path /path/to/pdf/files -monitoring true

Or you can provide these options in a configuration file (default is `config.json` in the current directory, but you can specify a different file with the `-config` flag):

>{ “APIKey”: “your-api-key”,
“Template”: “your-template”,
“ExtractionType”: “fast”,
“OutputJSON”: true,
“DirPath”: “/path/to/pdf/files”,
“Monitoring”: true }

And then run the application like this:

>go run main.go -config /path/to/config.json

## Dependencies

This application depends on several packages:

- `github.com/fsnotify/fsnotify`: For watching a directory for new files.
- `KlippaChallenge/pkg/api`: For making API requests.
- `KlippaChallenge/pkg/fileprocessing`: For processing files.
- `KlippaChallenge/pkg/monitor`: For monitoring a directory for new files.
- `KlippaChallenge/pkg/output`: For outputting results.

## Docker

You can also run this application in a Docker container. Here's how you can build and run the Docker image:

### Build

To build the Docker image, navigate to the directory containing the Dockerfile and run the following command:

>docker build -t klippa-ocr-challenge .

This command builds a Docker image using the Dockerfile in the current directory and tags it as `klippa-ocr-challenge`.

### Run

To run the Docker image, use the following command:

>docker run -v path\to\docs:/app/Docs -v path\to\outputs:/app/json-outputs klippa-app ./main -api-key="" -template= -monitoring=true -pdf-extraction=fast -output-json=true -dir-path="Docs"

This command runs the `klippa-ocr-challenge` Docker image. The `-v` flag is used to mount your local directories or files into the Docker container. Replace `/path/to/pdf/files` with the path to the directory containing your PDF files, and replace `path\to\outputs` with the desired outputs.

Please note that paths must be absolute. Also, if you're using Docker for Windows, you need to share your local drives with Docker.
