# Build Stage
FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Make the script executable
#RUN chmod +x start_server.sh

# Run the Bash script
CMD ["go", "run", "app/main.go"]