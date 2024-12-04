#!/bin/bash

echo "Starting script.sh with nohup..."
nohup go run app/main.go > app/app.log 2>&1 &
echo $! > app/app.pid
echo "script.sh is running in the background with PID: $(cat app/app.pid)"