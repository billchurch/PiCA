FROM golang:1.21-bullseye AS builder

WORKDIR /app
COPY . .

# Build the applications
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /pica ./cmd/pica
RUN CGO_ENABLED=0 GOOS=linux go build -o /pica-web ./cmd/pica-web

FROM debian:bullseye-slim

# Install dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    openssl \
    pcscd \
    pcsc-tools \
    yubikey-manager \
    yubico-piv-tool \
    && rm -rf /var/lib/apt/lists/*

# Copy built executables
COPY --from=builder /pica /usr/local/bin/
COPY --from=builder /pica-web /usr/local/bin/

# Copy web files and config
COPY web/html /app/web/html
COPY configs /app/configs

# Create required directories
RUN mkdir -p /app/certs /app/csrs /app/logs

WORKDIR /app

# Default command runs the web server
ENTRYPOINT ["/usr/local/bin/pica-web"]
CMD ["--config", "/app/configs/cfssl/sub-ca-config.json", "--cert", "/app/certs/sub-ca.pem", "--port", "8080", "--webroot", "/app/web/html", "--certdir", "/app/certs", "--csrdir", "/app/csrs"]

# Expose the API port
EXPOSE 8080
