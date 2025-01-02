package main

import (
	"testing"
)

func TestWildCardParsing(t *testing.T) {
	// Test the wildcard parsing
	// Test the wildcard parsing
	input := "192.168.*.*"
	expected := "192.168.0.0/16"
	cidrs, err := parseWildcard(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(cidrs) != 1 || cidrs[0].String() != expected {
		t.Fatalf("Expected %s, got %v", expected, cidrs)
	}
}

func TestInvalidInput(t *testing.T) {
	input := "300.300.300.300/24"
	_, err := parseCIDR(input)
	if err == nil {
		t.Fatalf("Expected error for invalid input, got none")
	}
}
