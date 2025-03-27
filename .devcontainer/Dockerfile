FROM mcr.microsoft.com/devcontainers/base:ubuntu-22.04

# Set environment variables
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin:/usr/local/go/bin
ENV DEBIAN_FRONTEND=noninteractive

# Install basic packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install YubiKey related packages - kept for compatibility with Linux hosts
# and for using YubiKey tools within the container
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    pcscd \
    pcsc-tools \
    libpcsclite1 \
    libpcsclite-dev \
    libccid \
    opensc \
    gnupg2 \
    openssl \
    socat \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create config directory for go tool
RUN mkdir -p /home/vscode/.config/go
RUN chown -R vscode:vscode /home/vscode/.config
