# --- Stage 1: Builder ---
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git for private module support
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /acs-proxy ./cmd/main.go

# --- Stage 2: Minimal runtime image ---
FROM gcr.io/distroless/static-debian11

# Set working directory (not strictly needed but good practice)
WORKDIR /

# Copy compiled binary
COPY --from=builder /acs-proxy /acs-proxy

# Copy config files
COPY --from=builder /app/internal/config /internal/config

# Copy swagger docs (JSON & YAML)
COPY --from=builder /app/docs /docs

# Expose ports
EXPOSE 3030
EXPOSE 9001

# Entrypoint
ENTRYPOINT ["/acs-proxy"]
