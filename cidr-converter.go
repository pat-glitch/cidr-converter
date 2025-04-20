// Enhanced CIDR Block Calculator with Expanded Input Formats in Go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
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

// deduplicateCIDRs removes duplicate CIDR blocks from the list.
func deduplicateCIDRs(cidrs []*net.IPNet) []*net.IPNet {
	seen := make(map[string]struct{})
	uniqueCIDRs := []*net.IPNet{}

	for _, cidr := range cidrs {
		cidrStr := cidr.String()
		if _, exists := seen[cidrStr]; !exists {
			seen[cidrStr] = struct{}{}
			uniqueCIDRs = append(uniqueCIDRs, cidr)
		}
	}
	return uniqueCIDRs
}

// ipBelongsToCIDR checks if the given IP belongs to any CIDR in the list.
func ipBelongsToCIDR(ipStr string, cidrs []*net.IPNet) ([]*net.IPNet, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	matchingCIDRs := []*net.IPNet{}
	for _, cidr := range cidrs {
		if cidr.Contains(ip) {
			matchingCIDRs = append(matchingCIDRs, cidr)
		}
	}
	return matchingCIDRs, nil
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

// mergeCIDRs merges a list of CIDR blocks into a minimal set.
func mergeCIDRs(cidrs []*net.IPNet) []*net.IPNet {
	sort.Slice(cidrs, func(i, j int) bool {
		return bytes.Compare(cidrs[i].IP, cidrs[j].IP) < 0
	})

	result := []*net.IPNet{}
	for _, cidr := range cidrs {
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
	sort.Slice(cidrs, func(i, j int) bool {
		return bytes.Compare(cidrs[i].IP, cidrs[j].IP) < 0
	})

	aggregated := []*net.IPNet{}
	for _, cidr := range cidrs {
		merged := false
		for i, agg := range aggregated {
			if canAggregate(agg, cidr) {
				aggregated[i] = mergeTwoCIDRs(agg, cidr)
				merged = true
				break
			}
		}
		if !merged {
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
	return bytes.Compare(a.IP, b.IP) == 0
}

// mergeTwoCIDRs merges two CIDR blocks into their parent CIDR.
func mergeTwoCIDRs(a, b *net.IPNet) *net.IPNet {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	ones, bits := a.Mask.Size()
	prefixLen := ones - 1
	parentIP := a.IP.Mask(net.CIDRMask(prefixLen, bits))
	return &net.IPNet{
		IP:   parentIP,
		Mask: net.CIDRMask(prefixLen, bits),
	}
}

// saveToJSON saves CIDRs to a JSON file.
func saveToJSON(filename string, cidrs []*net.IPNet) error {
	var cidrStrings []string
	for _, cidr := range cidrs {
		cidrStrings = append(cidrStrings, cidr.String())
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

	fmt.Println("Enter CIDR blocks, one per line. Enter an empty line to finish input:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		ipnet, err := parseCIDR(line)
		if err == nil {
			cidrs = append(cidrs, ipnet)
		} else {
			fmt.Printf("Invalid input: %s\n", err)
		}
	}

	// Deduplicate CIDRs
	cidrs = deduplicateCIDRs(cidrs)

	// Aggregate and merge CIDRs
	mergedCIDRs := aggregateCIDRs(mergeCIDRs(cidrs))

	fmt.Println("Merged and deduplicated CIDRs:")
	for _, cidr := range mergedCIDRs {
		fmt.Println(cidr)
	}

	// Check if an IP belongs to any CIDR
	fmt.Println("\nEnter an IP address to check:")
	if scanner.Scan() {
		ipInput := strings.TrimSpace(scanner.Text())
		matches, err := ipBelongsToCIDR(ipInput, mergedCIDRs)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else if len(matches) == 0 {
			fmt.Println("No matching CIDRs found.")
		} else {
			fmt.Println("Matching CIDRs:")
			for _, match := range matches {
				fmt.Println(match)
			}
		}
	}

	// Save merged CIDRs to a JSON file
	outputFile := "merged_cidrs.json"
	if err := saveToJSON(outputFile, mergedCIDRs); err != nil {
		fmt.Printf("Error saving JSON: %s\n", err)
	} else {
		fmt.Printf("\nMerged CIDRs saved to %s\n", outputFile)
	}
}
