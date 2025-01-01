CIDR Processing Tool
A command-line tool written in Go that processes IP address ranges in CIDR and wildcard notation. The tool can read input from standard input or CSV files and outputs merged CIDR blocks in JSON format.

Features

->Supports multiple input formats:
    1. CIDR notation (e.g., "192.168.1.0/24")
    2. Wildcard notation (e.g., "192.168.1.*")
    3. CSV files containing CIDR blocks

-> Validates IP ranges and CIDR blocks
-> Converts wildcard notation to CIDR format
-> Outputs results in JSON format
-> Comprehensive error handling