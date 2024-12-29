#!/bin/sh

printf "OpenAI API Key: "
read -s KEY

# Set environment variables
export D2D_KEY=$KEY

# Execute the doc2doc command with additional arguments
./doc2doc "$@"
