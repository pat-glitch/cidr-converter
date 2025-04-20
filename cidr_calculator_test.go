package main

import (
	"net"
	"os"
	"reflect"
	"testing"
)

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid CIDR",
			input:   "192.168.0.0/24",
			wantErr: false,
		},
		{
			name:    "Invalid CIDR format",
			input:   "192.168.0.0",
			wantErr: true,
		},
		{
			name:    "Invalid IP",
			input:   "256.256.256.256/24",
			wantErr: true,
		},
		{
			name:    "Invalid mask",
			input:   "192.168.0.0/33",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseCIDR(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCIDR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeduplicateCIDRs(t *testing.T) {
	_, net1, _ := net.ParseCIDR("192.168.0.0/24")
	_, net2, _ := net.ParseCIDR("192.168.1.0/24")
	_, net3, _ := net.ParseCIDR("192.168.0.0/24") // Duplicate of net1

	input := []*net.IPNet{net1, net2, net3}
	expected := []*net.IPNet{net1, net2}

	result := deduplicateCIDRs(input)
	if len(result) != len(expected) {
		t.Errorf("deduplicateCIDRs() returned %d CIDRs, expected %d", len(result), len(expected))
	}
}

func TestIPBelongsToCIDR(t *testing.T) {
	_, cidr1, _ := net.ParseCIDR("192.168.0.0/24")
	_, cidr2, _ := net.ParseCIDR("10.0.0.0/8")
	cidrs := []*net.IPNet{cidr1, cidr2}

	tests := []struct {
		name    string
		ip      string
		want    int
		wantErr bool
	}{
		{
			name:    "IP in first CIDR",
			ip:      "192.168.0.1",
			want:    1,
			wantErr: false,
		},
		{
			name:    "IP in second CIDR",
			ip:      "10.0.0.1",
			want:    1,
			wantErr: false,
		},
		{
			name:    "IP in no CIDR",
			ip:      "172.16.0.1",
			want:    0,
			wantErr: false,
		},
		{
			name:    "Invalid IP",
			ip:      "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, err := ipBelongsToCIDR(tt.ip, cidrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ipBelongsToCIDR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(matches) != tt.want {
				t.Errorf("ipBelongsToCIDR() returned %d matches, want %d", len(matches), tt.want)
			}
		})
	}
}

func TestParseWildcard(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid wildcard",
			input:   "192.168.*.*",
			want:    "192.168.0.0/16",
			wantErr: false,
		},
		{
			name:    "Invalid wildcard format",
			input:   "192.168.*",
			wantErr: true,
		},
		{
			name:    "Invalid IP parts",
			input:   "256.168.*.*",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWildcard(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWildcard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got[0].String() != tt.want {
				t.Errorf("parseWildcard() = %v, want %v", got[0].String(), tt.want)
			}
		})
	}
}

func TestMergeCIDRs(t *testing.T) {
	_, net1, _ := net.ParseCIDR("192.168.0.0/24")
	_, net2, _ := net.ParseCIDR("192.168.1.0/24")
	input := []*net.IPNet{net1, net2}

	result := mergeCIDRs(input)
	if len(result) != 2 {
		t.Errorf("mergeCIDRs() returned %d CIDRs, expected 2", len(result))
	}
}

func TestAggregateCIDRs(t *testing.T) {
	_, net1, _ := net.ParseCIDR("192.168.0.0/24")
	_, net2, _ := net.ParseCIDR("192.168.1.0/24")
	input := []*net.IPNet{net1, net2}

	result := aggregateCIDRs(input)
	if !reflect.DeepEqual(result, input) {
		t.Errorf("aggregateCIDRs() = %v, want %v", result, input)
	}
}

func TestSaveToJSON(t *testing.T) {
	_, net1, _ := net.ParseCIDR("192.168.0.0/24")
	_, net2, _ := net.ParseCIDR("192.168.1.0/24")
	cidrs := []*net.IPNet{net1, net2}

	tempFile := "test_cidrs.json"
	err := saveToJSON(tempFile, cidrs)
	if err != nil {
		t.Errorf("saveToJSON() error = %v", err)
	}
	// Clean up
	if err := deleteFile(tempFile); err != nil {
		t.Logf("Warning: Failed to delete test file: %v", err)
	}
}

// Helper function to delete test files
func deleteFile(filename string) error {
	return os.Remove(filename)
}
