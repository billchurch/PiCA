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

# Run the CLI application
run-cli:
	@echo "Running PiCA CLI..."
	./bin/$(BINARY_CLI)

# Run the web server
run-web:
	@echo "Running PiCA Web Server..."
	./bin/$(BINARY_WEB) \
		--config ./configs/cfssl/sub-ca-config.json \
		--cert ./certs/sub-ca.pem \
		--port 8080 \
		--webroot ./web/html \
		--certdir ./certs \
		--csrdir ./csrs

# Create directories
init:
	@echo "Creating required directories..."
	mkdir -p bin
	mkdir -p certs
	mkdir -p csrs
	mkdir -p logs

# Help
help:
	@echo "PiCA Makefile Help"
	@echo "------------------"
	@echo "make build       - Build both CLI and web applications"
	@echo "make build-cli   - Build only the CLI application"
	@echo "make build-web   - Build only the web server"
	@echo "make clean       - Clean up build artifacts"
	@echo "make test        - Run tests"
	@echo "make run-cli     - Run the CLI application"
	@echo "make run-web     - Run the web server"
	@echo "make init        - Create required directories"
	@echo "make build-root-ca-image - Build Root CA Raspberry Pi image"
	@echo "make build-sub-ca-image  - Build Sub CA Raspberry Pi image"
