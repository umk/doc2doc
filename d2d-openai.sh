#!/bin/sh

if [ -z "$D2D_KEY" ]; then
    printf "OpenAI API Key: "
    read -s KEY

    # Set environment variables
    export D2D_KEY=$KEY
fi

# Execute the doc2doc command with additional arguments
./doc2doc "$@"
