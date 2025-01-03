package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runCidrConverter(input string, inputFile string, outputFile string) (string, error) {
	cmd := exec.Command("go", "run", "cidr-converter.go", "-input", inputFile, "-output", outputFile)
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return out.String(), nil
}

func TestStandardInput(t *testing.T) {
	input := "192.168.0.0/24\n192.168.1.0/24\n"
	expected := `["192.168.0.0/23"]` //Expected output is a JSON format

	outputFile := "test_output.json"
	defer os.Remove(outputFile)

	_, err := runCidrConverter(input, "", outputFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify output file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Error reading output file: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestCSVInput(t *testing.T) {
	csvInput := `IP Range
	192.168.0.0-192.168.0.255
	192.168.1.0-192.168.1.255
	`

	inputFile := "test_input.csv"
	outputFile := "test_output.json"

	// Write CSV input to file
	err := os.WriteFile(inputFile, []byte(csvInput), 0644)
	if err != nil {
		t.Fatalf("Error writing CSV file: %v", err)
	}
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	_, err = runCidrConverter("", inputFile, outputFile)
	if err != nil {
		t.Errorf("Error running CIDR converter: %v", err)
	}

	//Verify output file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Error reading output file: %v", err)
	}

	expected := `["192.18.0.0/23"]`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestJSONInput(t *testing.T) {
	jsonInput := `{
		"cidrs": [
			"192.168.0.0/24",
			"192.168.1.0/24"
		]
	}`
	inputFile := "test_input.json"
	outputFile := "test_output.json"

	// Write JSON input to a file
	err := os.WriteFile(inputFile, []byte(jsonInput), 0644)
	if err != nil {
		t.Fatalf("Error creating test JSON file: %v", err)
	}
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	_, err = runCidrConverter("", inputFile, outputFile)
	if err != nil {
		t.Errorf("Error running CIDR converter: %v", err)
	}

	// Verify output
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Error reading output file: %v", err)
	}

	var result []string
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Errorf("Error parsing JSON output: %v", err)
	}

	expected := []string{"192.168.0.0/23"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestWildcardInput(t *testing.T) {
	input := "192.168.*.*\n"
	expected := `["192.168.0.0/16"]` // Expected merged output in JSON format

	outputFile := "test_output.json"
	defer os.Remove(outputFile)

	_, err := runCidrConverter(input, "", outputFile)
	if err != nil {
		t.Errorf("Error running CIDR converter: %v", err)
	}

	// Verify output file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Error reading output file: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected %s but got %s", expected, string(data))
	}
}
