#!/bin/bash

# Step 1: Build the binary
echo "Building the application..."
cd gf-emailparser
go build -o bin/gf-emailparser-binary ./cmd

# Check for successful build
if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi

# Check if a process is already running
if [ -f app.pid ]; then
    pid=$(cat app.pid)
    if ps -p $pid > /dev/null; then
        echo "Stopping the existing application (PID: $pid)..."
        kill $pid
        sleep 1
    fi
    rm app.pid
fi

# Start the server in the background and redirect output to a log file
echo "Starting the application..."
./bin/gf-emailparser-binary server.log 2>&1 & echo $! > app.pid
echo "Backend server started. Output is being logged to server.log."

# Allow some time for the server to start
sleep 2

exit 0