#!/bin/bash

# Chnage into Go directory
cd gf-emailparser

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