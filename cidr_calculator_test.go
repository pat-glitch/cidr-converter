package main

import (
	"net"
	"reflect"
	"testing"
)

// Helper function to parse CIDR blocks for test setup
func mustParseCIDR(cidr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return ipnet
}

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		input    string
		expected *net.IPNet
		hasError bool
	}{
		{"192.168.1.0/24", mustParseCIDR("192.168.1.0/24"), false},
		{"192.168.1.0/33", nil, true}, // Invalid CIDR
		{"invalid", nil, true},        // Non-CIDR input
	}

	for _, test := range tests {
		result, err := parseCIDR(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("parseCIDR(%q) error = %v, expected error: %v", test.input, err, test.hasError)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("parseCIDR(%q) = %v, expected: %v", test.input, result, test.expected)
		}
	}
}

func TestMergeCIDRs(t *testing.T) {
	input := []*net.IPNet{
		mustParseCIDR("192.168.1.0/24"),
		mustParseCIDR("192.168.2.0/24"),
		mustParseCIDR("192.168.1.128/25"),
	}
	expected := []*net.IPNet{
		mustParseCIDR("192.168.1.0/24"),
		mustParseCIDR("192.168.2.0/24"),
	}

	result := mergeCIDRs(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("mergeCIDRs(%v) = %v, expected: %v", input, result, expected)
	}
}

func TestAggregateCIDRs(t *testing.T) {
	input := []*net.IPNet{
		mustParseCIDR("192.168.1.0/24"),
		mustParseCIDR("192.168.2.0/24"),
	}
	expected := []*net.IPNet{
		mustParseCIDR("192.168.0.0/22"),
	}

	result := aggregateCIDRs(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("aggregateCIDRs(%v) = %v, expected: %v", input, result, expected)
	}
}

func TestParseWildcard(t *testing.T) {
	tests := []struct {
		input    string
		expected []*net.IPNet
		hasError bool
	}{
		{"192.168.*.*", []*net.IPNet{mustParseCIDR("192.168.0.0/16")}, false},
		{"192.168.1.*", []*net.IPNet{mustParseCIDR("192.168.1.0/24")}, false},
		{"invalid", nil, true},
	}

	for _, test := range tests {
		result, err := parseWildcard(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("parseWildcard(%q) error = %v, expected error: %v", test.input, err, test.hasError)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("parseWildcard(%q) = %v, expected: %v", test.input, result, test.expected)
		}
	}
}

func TestParseBinary(t *testing.T) {
	tests := []struct {
		input    string
		expected *net.IPNet
		hasError bool
	}{
		{"11000000101010000000000100000000/24", mustParseCIDR("192.168.1.0/24"), false},
		{"invalid", nil, true},
	}

	for _, test := range tests {
		result, err := parseBinary(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("parseBinary(%q) error = %v, expected error: %v", test.input, err, test.hasError)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("parseBinary(%q) = %v, expected: %v", test.input, result, test.expected)
		}
	}
}
