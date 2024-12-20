#!/bin/bash

# Load environment variables from .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found!"
    exit 1
fi

# Check if PORT_SERVER is set
if [ -z "$PORT_SERVER" ]; then
    echo "PORT_SERVER not set in .env"
    exit 1
fi

# Find and kill the process using PORT_SERVER
echo "Looking for processes using port $PORT_SERVER..."
PID=$(lsof -ti :$PORT_SERVER)

if [ -n "$PID" ]; then
    echo "Killing process using port $PORT_SERVER (PID: $PID)..."
    kill -9 $PID
    echo "Process $PID killed."
else
    echo "No process is using port $PORT_SERVER."
fi

# Start the server with nohup
echo "Starting script.sh with nohup..."
nohup go run app/main.go > app/app.log 2>&1 &
echo $! > app/app.pid
echo "script.sh is running in the background with PID: $(cat app/app.pid)"