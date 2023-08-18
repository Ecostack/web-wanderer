# Start with a base GoLang image
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the project files
COPY main.go ./

# Build the Go application
RUN go build -o app

# Layer 2: Final Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=builder /app/app .


# Set the entrypoint command to run the built executable
ENTRYPOINT ["./app"]