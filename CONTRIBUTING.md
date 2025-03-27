# Contributing to PiCA

Thank you for your interest in contributing to PiCA! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

1. Fork the repository
2. Clone your forked repository: `git clone https://github.com/billchurch/pica.git`
3. Create a new branch for your feature or bugfix: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Run tests: `make test`
6. Commit your changes: `git commit -m "Add your detailed commit message"`
7. Push to your branch: `git push origin feature/your-feature-name`
8. Open a Pull Request against the main repository

## Development Environment

### Prerequisites

- Go 1.21 or newer
- YubiKey with PIV support (for testing YubiKey functionality)
- PCSC daemon (pcscd)
- YubiKey tools (yubico-piv-tool, yubikey-manager)
- CFSSL toolkit

### Building and Testing

```bash
# Initialize directories
make init

# Build the applications
make build

# Run tests
make test

# Run the CLI application
make run-cli

# Run the web server
make run-web
```

## Project Structure

- `cmd/`: Command-line applications
- `internal/`: Internal packages
- `pkg/`: Public packages
- `web/`: Web interface
- `configs/`: Configuration files
- `rpi-images/`: Custom Raspberry Pi image configurations

## Pull Request Guidelines

- Follow Go best practices and code style
- Include tests for new features or bug fixes
- Keep changes focused on a single issue
- Update documentation as needed
- Provide a clear description of the changes in your PR

## Reporting Issues

- Use the GitHub issue tracker
- Provide detailed steps to reproduce
- Include information about your environment
- Mention any related issues or PRs

## Feature Requests

- Use the GitHub issue tracker
- Provide a clear description of the proposed feature
- Explain why the feature would be valuable
- Consider how the feature fits into the project's goals

## Security Issues

If you discover a security vulnerability, please do NOT open an issue. Email ...something...I'll figure it out in a bit... instead.

## Code Style

- Follow standard Go conventions
- Run `go fmt` and `go vet` before committing
- Use meaningful variable and function names
- Add comments for complex logic
- Write clear and concise commit messages

## Testing

- Write unit tests for new functionality
- Include integration tests for complex features
- Aim for high test coverage
- Test with real YubiKeys when possible

## Documentation

- Update documentation to reflect changes
- Use clear and concise language
- Include examples where appropriate
- Keep README.md and other docs up to date

## Licensing

By contributing to PiCA, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).

## Questions?

If you have any questions about contributing, feel free to open an issue asking for clarification.

Thank you for your contributions!
