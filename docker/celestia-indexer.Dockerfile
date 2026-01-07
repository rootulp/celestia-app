# Multi-stage build for celestia-indexer
# This Dockerfile clones and builds celestia-indexer from source

FROM golang:1.25.1-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git make

# Clone the celestia-indexer repository
WORKDIR /build
RUN git clone --depth 1 https://github.com/celenium-io/celestia-indexer.git .

# Download dependencies
RUN go mod download

# Build the application - try common entry points
RUN go build -o /build/celestia-indexer ./cmd/indexer 2>/dev/null || \
    go build -o /build/celestia-indexer ./cmd 2>/dev/null || \
    go build -o /build/celestia-indexer . 2>/dev/null || \
    (if [ -f Makefile ]; then make build && find . -name "*indexer*" -type f -executable -exec cp {} /build/celestia-indexer \; ; fi)

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS and curl for health checks
RUN apk --no-cache add ca-certificates curl bash

# Copy entire repository structure to /app - the indexer needs to run from the repo root
WORKDIR /app
COPY --from=builder /build /app

# The binary is already in /app from the COPY above
# Make it executable
RUN chmod +x /app/celestia-indexer

# Create symlinks or ensure the indexer can find the functions and views directories
# The indexer might be looking for "functions" and "views" in the current directory
RUN if [ -d /app/database/functions ]; then \
      ln -sf /app/database/functions /app/functions 2>/dev/null || true; \
    fi && \
    if [ -d /app/database/views ]; then \
      ln -sf /app/database/views /app/views 2>/dev/null || true; \
    fi

# Config file will be mounted via docker-compose volume

# Expose port (default for celestia-indexer API)
EXPOSE 9876

# Run the indexer from the repository root where all files are located
WORKDIR /app
CMD ["./celestia-indexer"]
