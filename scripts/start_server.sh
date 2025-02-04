#!/bin/bash

# Start the Python script
python3 scripts/live_playlist.py &
PYTHON_PID=$!

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    kill $PYTHON_PID
    exit 0
}

# Set trap for cleanup
trap cleanup SIGINT SIGTERM

# Run Go server
go run main.go

# If Go server exits, cleanup Python script
cleanup 