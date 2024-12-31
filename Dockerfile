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

# Build the application and the migration tool (if it's a Go file)
RUN CGO_ENABLED=0 GOOS=linux go build -o main . && go build -o migrate ./migrate/migrate.go

# Final stage (without Go installed)
FROM alpine:latest

# Add CA certificates (necessary for HTTPS connections)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Expose port 8080
EXPOSE 8080

# Command to run migrations and then start the application
CMD ["sh", "-c", "./migrate && ./main"]
