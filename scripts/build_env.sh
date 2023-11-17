#!/bin/sh

# Check if the path to the bash file is provided as an argument
if [ $# -eq 0 ]; then
  echo "Please provide the path to the bash file as an argument."
  exit 1
fi

# Extract the file name and directory path
file_path=$1
file_name=$(basename "$file_path")
directory_path=$(dirname "$file_path")

# Change to the directory containing the bash file
cd "$directory_path" || exit 1

# Database settings.
grep -E "export [A-Z_]+" "$file_name" | sed -E 's/export //' > .env
