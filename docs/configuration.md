# PiCA Configuration System

This document provides an overview of the PiCA configuration system, explaining how it works and how to use it effectively.

## Configuration Sources

PiCA supports multiple configuration sources with a clear priority order:

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Configuration files** (JSON or TOML)
4. **Default values** (lowest priority)

This allows for flexible deployment scenarios while maintaining consistency across different environments.

## Configuration Parameters

The following table shows the most common configuration parameters:

| Parameter Name    | Command Flag      | Environment Variable | Config File Key   | Default Value | Description                           |
|-------------------|-------------------|----------------------|-------------------|---------------|---------------------------------------|
| CA Type           | --ca-type         | CA_TYPE              | ca_type           | "sub"         | Type of CA: "root" or "sub"           |
| CA Config File    | --ca-config       | CA_CONFIG            | ca_config         |               | Path to CFSSL CA config JSON          |
| CA Certificate    | --ca-cert         | CA_CERT              | ca_cert           |               | Path to CA certificate                |
| Root CA Certificate | --root-ca-cert  | ROOT_CA_CERT         | root_ca_cert      |               | Path to Root CA certificate           |
| Crypto Provider   | --provider        | PICA_PROVIDER        | provider          |               | "yubikey" or "software"               |
| Key Slot          | --key-slot        | KEY_SLOT             | key_slot          | "82"          | YubiKey PIV slot (hex value)          |
| Web Port          | --port            | WEB_PORT             | web_port          | 8080          | Port for web server                   |
| Web Root          | --webroot         | WEB_ROOT             | web_root          | "./web/html"  | Directory for web UI files            |
| Certificate Dir   | --certdir         | CERT_DIR             | cert_dir          | "./certs"     | Directory for certificates            |
| CSR Dir           | --csrdir          | CSR_DIR              | csr_dir           | "./csrs"      | Directory for CSRs                    |
| Log Level         | --log-level       | LOG_LEVEL            | log_level         | "info"        | Logging level                         |

## Using Configuration Files

PiCA supports both JSON and TOML formats for configuration files.

### JSON Example

```json
{
  "ca_type": "sub",
  "ca_config": "./configs/cfssl/sub-ca-config.json",
  "ca_cert": "./certs/sub-ca.pem",
  "root_ca_cert": "./certs/root-ca.pem",
  "provider": "yubikey",
  "key_slot": "83",
  "web_port": 8443,
  "web_root": "./web/html",
  "cert_dir": "./certs",
  "csr_dir": "./csrs"
}
```

### TOML Example

```toml
# CA settings
ca_type = "sub"
ca_config = "./configs/cfssl/sub-ca-config.json"
ca_cert = "./certs/sub-ca.pem"
root_ca_cert = "./certs/root-ca.pem"

# Provider settings
provider = "yubikey"
key_slot = "83"

# Web server settings
web_port = 8443
web_root = "./web/html"

# Storage settings
cert_dir = "./certs"
csr_dir = "./csrs"
```

## Configuration Lookup

When no configuration file is explicitly provided, PiCA looks for configuration files in these locations (in order):

1. `./pica.json`
2. `./pica.toml`
3. `./configs/pica.json`
4. `./configs/pica.toml`
5. `$HOME/.pica/config.json`
6. `$HOME/.pica/config.toml`

## Using Environment Variables

Environment variables offer a convenient way to override configuration, especially in containerized environments:

```bash
# Basic configuration
export CA_TYPE=sub
export CA_CONFIG=./configs/cfssl/sub-ca-config.json
export CA_CERT=./certs/sub-ca.pem
export PICA_PROVIDER=yubikey
export WEB_PORT=8443

# Run PiCA
./bin/pica-web
```

## Using Command-Line Flags

Command-line flags provide the highest priority configuration method:

```bash
./bin/pica-web \
  --ca-type sub \
  --ca-config ./configs/cfssl/sub-ca-config.json \
  --ca-cert ./certs/sub-ca.pem \
  --port 8443 \
  --provider yubikey
```

## Default Configuration

You can create a default configuration file with:

```bash
make init-config
```

This creates a basic `pica.toml` file in the current directory that you can customize.

## Development vs. Production Settings

### Development Environment

For development without a physical YubiKey:

```bash
# Use software provider
export PICA_PROVIDER=software

# Run with sofware provider
./bin/pica-web --port 8080
```

Or use the convenience target:

```bash
make run-web-dev
```

### Production Environment

For production deployments:

```bash
# Use YubiKey provider
export PICA_PROVIDER=yubikey

# Run with hardware security
./bin/pica-web --port 443 --https --tls-cert ./certs/web.pem --tls-key ./certs/web.key
```

## Docker Environment

When running in Docker, environment variables are particularly useful:

```yaml
services:
  pica-sub-ca:
    image: pica-sub-ca
    environment:
      - CA_TYPE=sub
      - CA_CONFIG=/app/configs/cfssl/sub-ca-config.json
      - CA_CERT=/app/certs/sub-ca.pem
      - PICA_PROVIDER=yubikey
      - WEB_PORT=8443
    volumes:
      - ./certs:/app/certs
      - ./csrs:/app/csrs
```

## Technical Implementation

The configuration system is implemented in the `internal/config` package:

- `config.go`: Defines the Config struct and basic loading/parsing
- `flags.go`: Handles command-line flag registration and parsing
- `loader.go`: Orchestrates loading from all sources with proper priority

The Config struct is loaded early in application startup and then passed to the relevant components.

## Extending the Configuration

To add new configuration options:

1. Add the field to the `Config` struct in `internal/config/config.go` with appropriate tags
2. Use the new configuration option in your code
3. The loading mechanism will automatically handle the new option using existing logic

Example adding a new option:

```go
type Config struct {
    // Existing fields...
    
    // New option for certificate validity period
    CertValidityDays int `env:"CERT_VALIDITY_DAYS" flag:"cert-validity" config:"cert_validity_days" default:"365"`
}
```
