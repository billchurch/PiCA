# PiCA Configuration Examples

This directory contains example configuration files for PiCA in different formats.

## Configuration Files

- `root-ca.json` - Example configuration for a Root CA using JSON format
- `sub-ca.json` - Example configuration for a Sub CA using JSON format
- `pica.toml` - Example configuration using TOML format

## Usage

### Using a Configuration File

To use a configuration file, pass it to the `pica` or `pica-web` command:

```bash
# Start the CLI with a specific config file
./bin/pica --config ./configs/examples/root-ca.json

# Start the web server with a specific config file
./bin/pica-web --config ./configs/examples/sub-ca.json
```

### Environment Variables

You can also set configuration through environment variables:

```bash
# Configure as Root CA
export CA_TYPE=root
export CA_CONFIG=./configs/cfssl/root-ca-config.json
export CA_CERT=./certs/root-ca.pem
export KEY_SLOT=82
export PICA_PROVIDER=yubikey

# Run the CLI
./bin/pica
```

### Command-Line Options

Command-line options take precedence over environment variables and config files:

```bash
./bin/pica-web \
  --ca-type sub \
  --ca-config ./configs/cfssl/sub-ca-config.json \
  --ca-cert ./certs/sub-ca.pem \
  --root-ca-cert ./certs/root-ca.pem \
  --port 8443 \
  --webroot ./web/html \
  --certdir ./certs \
  --csrdir ./csrs
```

## Configuration Priority

PiCA loads configuration in the following order (from lowest to highest priority):

1. Default values
2. Configuration file (if specified)
3. Environment variables
4. Command-line arguments

Values from higher priority sources override those from lower priority sources.

## Default Locations

If no configuration file is specified, PiCA looks in these locations:

1. `./pica.json`
2. `./pica.toml`
3. `./configs/pica.json`
4. `./configs/pica.toml`
5. `$HOME/.pica/config.json`
6. `$HOME/.pica/config.toml`
