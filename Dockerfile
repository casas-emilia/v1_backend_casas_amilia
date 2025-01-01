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

# Build the application and the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o main . && CGO_ENABLED=0 GOOS=linux go build -o migrate ./migrate/migrate.go

# Final stage
FROM alpine:latest

# Add CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory
WORKDIR /root/

# Copy the binaries from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Expose port 8080
EXPOSE 8080

# Create a startup script
RUN echo '#!/bin/sh' > start.sh && \
    echo './migrate' >> start.sh && \
    echo './main' >> start.sh && \
    chmod +x start.sh

# Set the startup command
CMD ["./start.sh"]

