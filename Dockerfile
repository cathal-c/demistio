# Use the official Go image as the base image
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o demistio main.go

# Use a minimal base image for the final container
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/demistio .

# Expose the port your application runs on (if applicable)
EXPOSE 8080

# Command to run the binary
CMD ["./demistio"]
