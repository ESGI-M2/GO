#!/bin/bash

# Script to update all import statements from 'project/' to 'github.com/yourusername/go-orm/'
# Replace 'yourusername' with your actual GitHub username

echo "Updating import statements..."

# Find all .go files and replace the imports
find . -name "*.go" -type f -exec sed -i 's|project/|github.com/yourusername/go-orm/|g' {} \;

# Also update any references in test files
find . -name "*.go" -type f -exec sed -i 's|project/|github.com/yourusername/go-orm/|g' {} \;

echo "Import statements updated!"
echo "Please replace 'yourusername' with your actual GitHub username in the updated files." 