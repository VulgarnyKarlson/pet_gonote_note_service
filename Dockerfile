# Start from a small, secure base image
FROM golang:1.21.0-alpine AS builder
RUN apk update
RUN apk add --no-cache make
ENV GO111MODULE=on
WORKDIR /go/src/app
RUN go env -w GOPRIVATE=gitlab.karlson.dev

# Copy the Go module files
COPY go.* ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN make build-linux

# Create a minimal production image
FROM alpine:latest

# It's essential to regularly update the packages within the image to include security patches
RUN apk update && apk upgrade

# Reduce image size
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Avoid running code as a root user
RUN adduser -D appuser
USER appuser

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /go/src/app/bin/main .

# Set any environment variables required by the application
ENV HTTP_ADDR=:8080

# Expose the port that the application listens on
EXPOSE 8080

# Run the binary when the container starts
CMD ["./main"]
