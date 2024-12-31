# Build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# List the contents of /app to verify the files
RUN ls -l /app

# Build the application with specific flags for a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Add CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the migration folder
COPY --from=builder /app/migrate ./migrate

# Expose port 8080
EXPOSE 8080

# Command to run migrations and then start the application
CMD ["sh", "-c", "./main"]
