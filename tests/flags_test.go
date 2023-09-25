package tests

import (
	"flag"
	"os"
	"testing"
)

func TestFlags(t *testing.T) {
	// Capture the original command-line arguments
	originalArgs := os.Args

	// Create a new set of command-line arguments for testing
	testArgs := []string{"-api-key", "my-api-key", "-template", "my-template"}

	// Replace the command-line arguments for testing
	os.Args = testArgs

	// Parse the new arguments
	flag.Parse()

	// Access flag values and assert their correctness
	apiKey := *flag.String("api-key", "", "Klippa OCR API key")
	template := *flag.String("template", "", "OCR template")

	if apiKey != "my-api-key" {
		t.Errorf("Expected API key 'my-api-key', but got '%s'", apiKey)
	}
	if template != "my-template" {
		t.Errorf("Expected template 'my-template', but got '%s'", template)
	}

	// Reset the command-line arguments to their original values
	os.Args = originalArgs
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
