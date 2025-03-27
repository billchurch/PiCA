# PiCA Development Guide

This document provides instructions for setting up and using the development environment for PiCA using VS Code's devcontainer feature.

## Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [Visual Studio Code](https://code.visualstudio.com/)
- [Remote - Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) for VS Code

For YubiKey testing (optional):
- A YubiKey device
- Host system with YubiKey support (drivers installed)

## Setting Up the Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/PiCA.git
   cd PiCA
   ```

2. Open the project in VS Code:
   ```bash
   code .
   ```

3. VS Code will detect the devcontainer configuration and prompt you to reopen the project in a container. Click "Reopen in Container".

   Alternatively, you can click the green button in the lower-left corner of VS Code and select "Reopen in Container".

4. The devcontainer will build and start automatically. This may take a few minutes the first time.

## Development Container Features

The devcontainer includes:

- Go development environment
- YubiKey support (pcscd, yubikey-manager, etc.)
- Docker-outside-of-docker support
- Git configuration
- VS Code extensions for Go development
- CFSSL toolkit
- Various Go development tools

## Working with YubiKey

### Testing YubiKey Access

1. Insert your YubiKey into a USB port on your host machine.

2. Run the "Check YubiKey Connection" task:
   - Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on macOS)
   - Type "Tasks: Run Task"
   - Select "Check YubiKey Connection"

   You should see information about your YubiKey if it's properly connected.

3. If the YubiKey is not detected, try running the "Start PCSC Daemon" task.

### Running the Provider Example

You can run the provider example to test both the software and YubiKey providers:

1. Auto-detection mode (will use YubiKey if available, otherwise software):
   - Run the "Run Provider Example (Auto)" task

2. Force software provider (for testing without YubiKey):
   - Run the "Run Provider Example (Software)" task

3. Force YubiKey provider (will fail if no YubiKey is connected):
   - Run the "Run Provider Example (YubiKey)" task

## Debugging

VS Code is configured for debugging Go applications. To start debugging:

1. Open the file you want to debug (e.g., `examples/provider_example.go`).

2. Set breakpoints by clicking in the gutter.

3. Press F5 or select "Run > Start Debugging".

4. Choose the appropriate launch configuration from the dropdown menu:
   - "Launch Provider Example" (auto-detection)
   - "Launch Provider Example (Software)"
   - "Launch Provider Example (YubiKey)"
   - "Launch PiCA CLI"
   - "Launch PiCA Web"

## Running Tests

To run tests:

1. From the VS Code menu, select "Terminal > Run Task..." and choose "Test PiCA".

   Alternatively, you can use the test explorer in VS Code's sidebar or run individual tests by clicking the "Run Test" link above each test function.

## Building

To build the application:

1. From the VS Code menu, select "Terminal > Run Task..." and choose "Build PiCA".

## VS Code Tasks

The following tasks are available:

- **Build PiCA**: Builds the entire project
- **Test PiCA**: Runs all tests
- **Run Provider Example (Auto)**: Runs with auto-detection
- **Run Provider Example (Software)**: Forces software provider
- **Run Provider Example (YubiKey)**: Forces YubiKey provider
- **Check YubiKey Connection**: Tests YubiKey connectivity
- **Start PCSC Daemon**: Starts the PCSC daemon for YubiKey
- **Generate Go Mocks**: Generates mock implementations
- **Lint Code**: Runs the linter

## Environment Variables

- `PICA_PROVIDER`: Controls provider selection
  - `software`: Forces software provider
  - `yubikey`: Forces YubiKey provider
  - Not set: Auto-detection

## Crypto Provider Options

PiCA supports different cryptographic provider types controlled via the `PICA_PROVIDER` environment variable:

- `export PICA_PROVIDER=yubikey` - Force YubiKey provider
- `export PICA_PROVIDER=software` - Force software provider
- Not set - Auto-detect (YubiKey if available, otherwise software)

The software provider stores keys and certificates in:
- `~/.pica/keys/` - Private key storage (protected by file permissions)
- `~/.pica/certs/` - Certificate storage

When testing without a YubiKey, use:

```bash
export PICA_PROVIDER=software
```

## USB Passthrough

USB devices (like YubiKey) are passed through to the container using the mount configuration in `devcontainer.json`. This should work automatically, but if you're having issues:

1. Make sure your host system has the necessary YubiKey drivers installed.
2. Try unplugging and re-plugging the YubiKey after the container is running.
3. Check the PCSC daemon status with `sudo systemctl status pcscd`.

## Troubleshooting

### YubiKey Not Detected

1. Check if the YubiKey is detected on your host system:
   ```bash
   # On macOS
   system_profiler SPUSBDataType | grep -A 10 Yubico
   
   # On Linux
   lsusb | grep Yubico
   ```

2. Make sure the PCSC daemon is running inside the container:
   ```bash
   sudo systemctl status pcscd
   ```

3. If it's not running, start it:
   ```bash
   sudo systemctl start pcscd
   ```

### Go Modules Issues

If you're having issues with Go modules:

1. Make sure your Go version is correct:
   ```bash
   go version
   ```

2. Try cleaning the module cache:
   ```bash
   go clean -modcache
   ```

3. Update dependencies:
   ```bash
   go get -u ./...
   ```
