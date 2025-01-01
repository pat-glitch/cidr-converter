package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// parseCIDR validates and returns a CIDR block
func parseCIDR(input string) (*net.IPNet, error) {
	ip, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return nil, fmt.Errorf("Invalid CIDR block: %v", input)
	}
	if ip == nil || ipnet == nil {
		return nil, fmt.Errorf("Invalid IP or CIDR range")
	}
	return ipnet, nil
}

// mergeCIDRs merges a list of CIDR blocks into a single CIDR block
func mergeCIDRs(cidrs []*net.IPNet) (*net.IPNet, error) {
	// This basic implementation does not yet collapse CIDRs.
	// Future versions will collapse CIDRs.
	// For now, it just returns the first CIDR block.
	if len(cidrs) == 0 {
		return nil, fmt.Errorf("No CIDR blocks to merge")
	}
	return cidrs[0], nil
}

// Main function
func main() {
	var cidrs []*net.IPNet

	// Read input from stdin
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter CIDR blocks (one per line).Press Ctrl+D(Linux/Mac) or Ctrl+Z(Windows) to end input:")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		ipnet, err := parseCIDR(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		cidrs = append(cidrs, ipnet)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
		result, err := mergeCIDRs(cidrs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error merging CIDRs: %v\n", err)
			os.Exit(1)
		}

		// Print the merged CIDR
		fmt.Println("Merged CIDR block:")
		fmt.Println(result.String())
	}
}
