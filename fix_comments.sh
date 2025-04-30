#!/bin/bash

# Find all Go files with the pattern "// -" and replace with "//   -"
find . -name "*.go" -type f -exec grep -l "// -" {} \; | while read file; do
  echo "Processing $file..."
  # Replace "// -" with "//   -" in the file
  sed -i 's|// -|//   -|g' "$file"
done

echo "All files processed."
