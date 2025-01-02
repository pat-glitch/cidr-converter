#!/bin/bash

echo "Running tests..."

for input_file in tests/input*.txt tests/input*.csv; do
  base_name=$(basename "$input_file" .txt)
  base_name=$(basename "$base_name" .csv)
  expected_output="tests/${base_name}.json"

  echo "Testing $input_file..."
  ./cidr_calculator "$input_file" > "tests/actual_${base_name}.json"

  if diff "tests/actual_${base_name}.json" "$expected_output"; then
    echo "✅ $input_file passed!"
  else
    echo "❌ $input_file failed! See tests/actual_${base_name}.json"
  fi
done
