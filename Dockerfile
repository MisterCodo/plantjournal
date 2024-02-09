# First stage: Build the Go application
FROM golang:1.22.0-bookworm AS builder

WORKDIR /app

# Cache go.mod
COPY go.mod go.sum ./
RUN go mod download

# Copy remainder
COPY . .

# Build the Go application with CGO enabled, needed for sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/plantjournal /app/cmd/plantjournal/main.go

# Second stage: Launch the Go application
FROM debian:12.4-slim

# Create the necessary directories
WORKDIR /home/appuser
RUN mkdir -p /home/appuser/app

# Copy the Go application from the builder stage
COPY --from=builder /app/plantjournal /home/appuser/app/plantjournal

# Copy the HTML template files from the builder stage
COPY --from=builder /app/templates/ /home/appuser/templates/

# Copy static files from the builder stage
COPY --from=builder /app/static/ /home/appuser/static/

# Set the entrypoint for the application
ENTRYPOINT ["/home/appuser/app/plantjournal"]
