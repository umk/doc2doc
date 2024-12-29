#!/bin/sh

# Check if the "model" parameter is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <model> [additional arguments...]"
  exit 1
fi

# Set environment variables
export D2D_BASE_URL=http://localhost:11434/v1/
export D2D_KEY=1
export D2D_MODEL=$1

# Execute the doc2doc command with additional arguments
./doc2doc "${@:2}"
