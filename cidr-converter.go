// Enhanced CIDR Block Calculator with Expanded Input Formats in Go
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	// Filter out nil entries
	validCIDRs := []*net.IPNet{}
	for _, cidr := range cidrs {
		if cidr != nil {
			validCIDRs = append(validCIDRs, cidr)
		}
	}

	sort.Slice(validCIDRs, func(i, j int) bool {
		return bytes.Compare(validCIDRs[i].IP, validCIDRs[j].IP) < 0
	})

	result := []*net.IPNet{}
	for _, cidr := range validCIDRs {
		if len(result) == 0 {
			result = append(result, cidr)
			continue
		}
		last := result[len(result)-1]
		if last.Contains(cidr.IP) {
			continue
		}
		result = append(result, cidr)
	}
	return result
}

// aggregateCIDRs aggregates smaller subnets into larger ones when possible.
func aggregateCIDRs(cidrs []*net.IPNet) []*net.IPNet {
	// Sort CIDRs to make aggregation easier
	sort.Slice(cidrs, func(i, j int) bool {
		return bytes.Compare(cidrs[i].IP, cidrs[j].IP) < 0
	})

	aggregated := []*net.IPNet{}
	for _, cidr := range cidrs {
		merged := false
		for i, agg := range aggregated {
			// Check if the CIDR can be merged with the current aggregated CIDR
			if canAggregate(agg, cidr) {
				// Merge and update the aggregated CIDR
				aggregated[i] = mergeTwoCIDRs(agg, cidr)
				merged = true
				break
			}
		}
		if !merged {
			// If no merge happened, just append the current CIDR
			aggregated = append(aggregated, cidr)
		}
	}
	return aggregated
}

// canAggregate checks if two CIDR blocks can be aggregated into a larger block.
func canAggregate(a, b *net.IPNet) bool {
	if a == nil || b == nil {
		return false
	}
	onesA, bitsA := a.Mask.Size()
	onesB, bitsB := b.Mask.Size()
	if bitsA != bitsB || onesA != onesB {
		return false
	}
	// Check if the two CIDRs are adjacent
	diff := bytes.Compare(a.IP, b.IP)
	return diff == 1 || diff == -1
}

// mergeTwoCIDRs merges two CIDR blocks into their parent CIDR.
func mergeTwoCIDRs(a, b *net.IPNet) *net.IPNet {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	ones, _ := a.Mask.Size()
	prefixLen := ones - 1
	parentIP := a.IP.Mask(net.CIDRMask(prefixLen, 32))
	return &net.IPNet{
		IP:   parentIP,
		Mask: net.CIDRMask(prefixLen, 32),
	}
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

// parseBinary converts binary string representations to CIDR blocks.
func parseBinary(input string) (*net.IPNet, error) {
	binaryRegex := regexp.MustCompile(`^[01]{32}/\d{1,2}$`)
	if !binaryRegex.MatchString(input) {
		return nil, fmt.Errorf("invalid binary CIDR notation: %s", input)
	}
	parts := strings.Split(input, "/")
	binaryIP := parts[0]
	prefix, err := strconv.Atoi(parts[1])
	if err != nil || prefix < 0 || prefix > 32 {
		return nil, fmt.Errorf("invalid prefix length: %s", parts[1])
	}
	ip := net.IP{0, 0, 0, 0}
	for i := 0; i < 32; i++ {
		if binaryIP[i] == '1' {
			ip[i/8] |= (1 << uint(7-i%8))
		}
	}
	cidr := fmt.Sprintf("%s/%d", ip.String(), prefix)
	return parseCIDR(cidr)
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

func parseJSON(filename string) ([]*net.IPNet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var cidrStrings []string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cidrStrings); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	var cidrs []*net.IPNet
	for _, entry := range cidrStrings {
		ipnet, err := parseCIDR(entry)
		if err == nil {
			cidrs = append(cidrs, ipnet)
		}
	}
	return cidrs, nil
}

func saveToJSON(filename string, cidrs []*net.IPNet) error {
	var cidrStrings []string
	for _, cidr := range cidrs {
		if cidr != nil {
			cidrStrings = append(cidrStrings, cidr.String())
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cidrStrings); err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil
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

		var wg sync.WaitGroup
		inputChan := make(chan string)
		outputChan := make(chan *net.IPNet, 100)

		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range inputChan {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				var ipnet *net.IPNet
				var err error
				if strings.Contains(line, "*") {
					wildcardCidrs, err := parseWildcard(line)
					if err == nil {
						for _, wc := range wildcardCidrs {
							outputChan <- wc
						}
					}
				} else if strings.Contains(line, "0") || strings.Contains(line, "1") {
					ipnet, err = parseBinary(line)
				} else {
					ipnet, err = parseCIDR(line)
				}

				if err == nil && ipnet != nil {
					outputChan <- ipnet
				} else {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			}
		}()

		go func() {
			for scanner.Scan() {
				inputChan <- scanner.Text()
			}
			close(inputChan)
		}()

		go func() {
			wg.Wait()
			close(outputChan)
		}()

		for ipnet := range outputChan {
			if ipnet != nil {
				cidrs = append(cidrs, ipnet)
			}
		}
	} else if strings.HasSuffix(inputType, ".csv") {
		// Parse from CSV file
		csvCidrs, err := parseCSV(inputType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
			os.Exit(1)
		}
		cidrs = append(cidrs, csvCidrs...)
	} else if strings.HasSuffix(inputType, ".json") {
		// Parse from JSON file
		jsonCidrs, err := parseJSON(inputType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading JSON: %v\n", err)
			os.Exit(1)
		}
		cidrs = append(cidrs, jsonCidrs...)
	} else {
		fmt.Fprintf(os.Stderr, "Unsupported input format\n")
		os.Exit(1)
	}
	// Merge CIDRs(cidrs)
	result := aggregateCIDRs(mergeCIDRs(cidrs))

	// Output merged CIDR blocks in JSON format
	outputFile := "merged_cidrs.json"
	if err := saveToJSON(outputFile, result); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Merged CIDR blocks saved to %s\n", outputFile)
}
