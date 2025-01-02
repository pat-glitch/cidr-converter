// Enhanced CIDR Block Calculator with Expanded Input Formats in Go
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// parseCIDR validates and returns a CIDR block.
func parseCIDR(input string) (*net.IPNet, error) {
	ip, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR block: %v", input)
	}
	if ip == nil || ipnet == nil {
		return nil, fmt.Errorf("invalid IP or CIDR range")
	}
	return ipnet, nil
}

// mergeCIDRs merges a list of CIDR blocks into a minimal set.
func mergeCIDRs(cidrs []*net.IPNet) []*net.IPNet {
	// This basic implementation does not yet collapse CIDRs.
	// Future iterations can add functionality for collapsing adjacent ranges.
	return cidrs
}

// parseWildcard converts wildcard notation (e.g., 192.168.*.*) to CIDR blocks.
func parseWildcard(input string) ([]*net.IPNet, error) {
	wildcardRegex := regexp.MustCompile(`^((?:\d{1,3}|\*)\.){3}(?:\d{1,3}|\*)$`)
	if !wildcardRegex.MatchString(input) {
		return nil, fmt.Errorf("invalid wildcard notation: %s", input)
	}

	octets := strings.Split(input, ".")
	var ipRange string
	for _, octet := range octets {
		if octet == "*" {
			ipRange += "0."
		} else {
			ipRange += octet + "."
		}
	}
	ipRange = strings.TrimSuffix(ipRange, ".") + "/"

	prefix := 32
	for _, octet := range octets {
		if octet == "*" {
			prefix -= 8
		}
	}

	cidr := fmt.Sprintf("%s%d", ipRange, prefix)
	ipnet, err := parseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return []*net.IPNet{ipnet}, nil
}

func parseCSV(filename string) ([]*net.IPNet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var cidrs []*net.IPNet

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	for _, record := range records {
		if len(record) < 1 {
			continue
		}
		entry := strings.TrimSpace(record[0])
		ipnet, err := parseCIDR(entry)
		if err == nil {
			cidrs = append(cidrs, ipnet)
		}
	}
	return cidrs, nil
}

func main() {
	var cidrs []*net.IPNet

	inputType := "stdin"
	if len(os.Args) > 1 {
		inputType = os.Args[1]
	}

	if inputType == "stdin" {
		// Read input from stdin
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter CIDR blocks, one per line. Press Ctrl+D (Linux/Mac) or Ctrl+Z (Windows) to end input:")

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			if strings.Contains(line, "*") {
				wildcardCidrs, err := parseWildcard(line)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					continue
				}
				cidrs = append(cidrs, wildcardCidrs...)
			} else {
				ipnet, err := parseCIDR(line)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					continue
				}
				cidrs = append(cidrs, ipnet)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	} else if strings.HasSuffix(inputType, ".csv") {
		// Parse from CSV file
		csvCidrs, err := parseCSV(inputType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
			os.Exit(1)
		}
		cidrs = append(cidrs, csvCidrs...)
	} else {
		fmt.Fprintf(os.Stderr, "Unsupported input format\n")
		os.Exit(1)
	}

	result := mergeCIDRs(cidrs)

	fmt.Println("Merged CIDR blocks in JSON:")
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
