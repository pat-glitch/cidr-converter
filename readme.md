# CIDR Convert

A command-line utility written in Go that processes, validates, and merges IP address ranges in various formats. The tool supports CIDR notation, wildcard notation, and multiple input/output formats.

## Features

### Input Processing
- Multiple input formats supported:
  - CIDR notation (e.g., "192.168.1.0/24")
  - Wildcard notation (e.g., "192.168.1.*")
  - CSV files containing CIDR blocks
  - JSON files containing CIDR blocks
- Interactive stdin mode for manual input

### CIDR Operations
- Validates IP ranges and CIDR blocks
- Converts wildcard notation to CIDR format
- Merges overlapping CIDR blocks
- Sorts CIDR blocks for optimal organization

### Output Handling
- Automatically saves merged results to JSON file
- Pretty-printed JSON output
- Comprehensive error handling and reporting

## Installation

Ensure you have Go installed on your system, then:

```bash
git clone [repository-url]
cd [repository-name]
go build
```

## Usage

The tool supports three input modes:

### 1. Standard Input Mode

```bash
./cidr-processor
# Enter CIDR blocks interactively, one per line:
192.168.1.0/24
10.0.0.*
# Press Ctrl+D (Linux/Mac) or Ctrl+Z (Windows) to end input
```

### 2. CSV File Mode

```bash
./cidr-processor input.csv
```

CSV file format:
```csv
192.168.1.0/24
10.0.0.0/8
172.16.0.0/12
```

### 3. JSON File Mode

```bash
./cidr-processor input.json
```

JSON file format:
```json
[
  "192.168.1.0/24",
  "10.0.0.0/8",
  "172.16.0.0/12"
]
```

## Output

The tool saves merged CIDR blocks to `test_output.json`:

```json
[
  "10.0.0.0/8",
  "172.16.0.0/12",
  "192.168.1.0/24"
]
```

## Error Handling

The tool handles various error cases:
- Invalid CIDR blocks
- Malformed wildcard notation
- File reading/writing errors
- CSV/JSON parsing issues

## Requirements

- Go 1.x or higher
- No external dependencies (uses only Go standard library)

Inspired from Andy Walker's [cidr-convert repo](https://github.com/flowchartsman/cidr-convert)