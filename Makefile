.PHONY: build build-cli build-web clean test

BINARY_CLI=pica
BINARY_WEB=pica-web

all: clean build

build: build-cli build-web

build-cli:
	@echo "Building PiCA CLI..."
	go build -o bin/$(BINARY_CLI) ./cmd/pica

build-web:
	@echo "Building PiCA Web Server..."
	go build -o bin/$(BINARY_WEB) ./cmd/pica-web

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	mkdir -p bin/

# Root CA image
build-root-ca-image:
	@echo "Building Root CA image using rpi-image-gen..."
	cd rpi-image-gen && ./build.sh -c $(PWD)/rpi-images/root-ca/config/config.ini

# Sub CA image
build-sub-ca-image:
	@echo "Building Sub CA image using rpi-image-gen..."
	cd rpi-image-gen && ./build.sh -c $(PWD)/rpi-images/sub-ca/config/config.ini

# Install to YubiKey
install-to-yubikey:
	@echo "Installing to YubiKey..."
	@echo "This feature is not yet implemented"

# Run the CLI application with different configurations
run-cli:
	@echo "Running PiCA CLI with default configuration..."
	./bin/$(BINARY_CLI)

run-cli-root:
	@echo "Running PiCA CLI as Root CA..."
	./bin/$(BINARY_CLI) --config ./configs/examples/root-ca.json

run-cli-sub:
	@echo "Running PiCA CLI as Sub CA..."
	./bin/$(BINARY_CLI) --config ./configs/examples/sub-ca.json

# Run the web server with different configurations
run-web:
	@echo "Running PiCA Web Server with default configuration..."
	./bin/$(BINARY_WEB)

run-web-config:
	@echo "Running PiCA Web Server with example configuration..."
	./bin/$(BINARY_WEB) --config ./configs/examples/sub-ca.json

# Run with software provider (for development without YubiKey)
run-web-dev:
	@echo "Running PiCA Web Server in development mode (software provider)..."
	PICA_PROVIDER=software ./bin/$(BINARY_WEB) --config ./configs/examples/sub-ca.json

# Run configuration example
run-config-example:
	@echo "Running configuration example..."
	go run examples/config/main.go

# Configuration examples with different sources
run-config-env:
	@echo "Running config example with environment variables..."
	CA_TYPE=root PICA_PROVIDER=software go run examples/config/main.go

run-config-args:
	@echo "Running config example with command-line args..."
	go run examples/config/main.go --ca-type=root --provider=yubikey --port=8443

run-config-file:
	@echo "Running config example with configuration file..."
	go run examples/config/main.go --config=./configs/examples/root-ca.json

# Create directories
init:
	@echo "Creating required directories..."
	mkdir -p bin
	mkdir -p certs
	mkdir -p csrs
	mkdir -p logs
	mkdir -p db

# Create basic configuration file
init-config:
	@echo "Creating basic configuration files..."
	mkdir -p configs
	cp -n configs/examples/pica.toml ./pica.toml || true
	@echo "Configuration file created: pica.toml"

# Help
help:
	@echo "PiCA Makefile Help"
	@echo "------------------"
	@echo "make build       - Build both CLI and web applications"
	@echo "make build-cli   - Build only the CLI application"
	@echo "make build-web   - Build only the web server"
	@echo "make clean       - Clean up build artifacts"
	@echo "make test        - Run tests"
	@echo "make init        - Create required directories"
	@echo "make init-config - Create basic configuration file"
	@echo
	@echo "Running applications:"
	@echo "make run-cli       - Run the CLI application"
	@echo "make run-cli-root  - Run the CLI application as Root CA"
	@echo "make run-cli-sub   - Run the CLI application as Sub CA"
	@echo "make run-web       - Run the web server"
	@echo "make run-web-config - Run the web server with example configuration"
	@echo "make run-web-dev   - Run the web server with software provider"
	@echo
	@echo "Configuration examples:"
	@echo "make run-config-example - Run basic configuration example"
	@echo "make run-config-env     - Run example with environment variables"
	@echo "make run-config-args    - Run example with command-line args"
	@echo "make run-config-file    - Run example with configuration file"
	@echo
	@echo "Raspberry Pi images:"
	@echo "make build-root-ca-image - Build Root CA Raspberry Pi image"
	@echo "make build-sub-ca-image  - Build Sub CA Raspberry Pi image"
