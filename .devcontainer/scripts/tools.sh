#!/bin/bash
set -e

echo "Setting up development environment for PiCA..."

# Configure git
git config --global --add safe.directory ${PWD}

echo "Installing Go tools..."
# Install essential Go tools
# go install golang.org/x/tools/cmd/goimports@latest
# go install github.com/go-delve/delve/cmd/dlv@latest

# Install CFSSL toolkit
echo "Installing CFSSL toolkit..."
go install github.com/cloudflare/cfssl/cmd/cfssl@latest
go install github.com/cloudflare/cfssl/cmd/cfssljson@latest

echo "Setup complete. Development environment ready."
