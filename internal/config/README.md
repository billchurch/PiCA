# PiCA Configuration System

This package provides a unified configuration system for PiCA that supports loading settings from three sources, with the following priority order:

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

## Usage

### Basic Usage

```go
import "github.com/billchurch/PiCA/internal/config"

func main() {
    // Load configuration with command-line arguments
    cfg, err := config.Load(os.Args[1:], "")
    if err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }

    // Use the configuration
    fmt.Printf("Web server port: %d\n", cfg.WebPort)
    fmt.Printf("Certificate directory: %s\n", cfg.CertDir)
}
```

### Specifying a Configuration File

```go
// Load with a specific config file
cfg, err := config.Load(os.Args[1:], "/path/to/config.json")
```

### Setting Configuration Options

#### Via Command-line

```bash
./pica --port 8443 --ca-cert ./certs/sub-ca.pem --provider yubikey
```

#### Via Environment Variables

```bash
export WEB_PORT=8443
export CA_CERT=./certs/sub-ca.pem
export PICA_PROVIDER=yubikey
./pica
```

#### Via Configuration File (JSON)

```json
{
  "web_port": 8443,
  "ca_cert": "./certs/sub-ca.pem",
  "provider": "yubikey"
}
```

#### Via Configuration File (TOML)

```toml
web_port = 8443
ca_cert = "./certs/sub-ca.pem"
provider = "yubikey"
```

## Configuration Options

| Option Name       | Environment Variable | Flag            | Config Key     | Default  | Description                             |
|-------------------|----------------------|-----------------|----------------|----------|-----------------------------------------|
| LogLevel          | LOG_LEVEL            | --log-level     | log_level      | info     | Logging level                           |
| ConfigDir         | CONFIG_DIR           | --config-dir    | config_dir     | ./configs| Configuration directory                 |
| CAType            | CA_TYPE              | --ca-type       | ca_type        | sub      | CA type (root or sub)                   |
| CAConfigFile      | CA_CONFIG            | --ca-config     | ca_config      |          | Path to CA config file                  |
| CACertFile        | CA_CERT              | --ca-cert       | ca_cert        |          | Path to CA certificate file             |
| CRLFile           | CRL_FILE             | --crl-file      | crl_file       |          | Path to CRL file                        |
| RootCACertFile    | ROOT_CA_CERT         | --root-ca-cert  | root_ca_cert   |          | Path to Root CA certificate file        |
| RootCAConfigFile  | ROOT_CA_CONFIG       | --root-ca-config| root_ca_config |          | Path to Root CA config file             |
| CAProfile         | CA_PROFILE           | --ca-profile    | ca_profile     |          | CA profile to use                       |
| ProviderType      | PICA_PROVIDER        | --provider      | provider       |          | Crypto provider type (yubikey/software) |
| KeySlot           | KEY_SLOT             | --key-slot      | key_slot       | 82       | YubiKey PIV slot to use                 |
| WebPort           | WEB_PORT             | --port          | web_port       | 8080     | Web server port                         |
| WebRoot           | WEB_ROOT             | --webroot       | web_root       | ./web/html | Web server root directory              |
| WebTLSCert        | WEB_TLS_CERT         | --tls-cert      | web_tls_cert   |          | TLS certificate for HTTPS               |
| WebTLSKey         | WEB_TLS_KEY          | --tls-key       | web_tls_key    |          | TLS key for HTTPS                       |
| EnableHTTPS       | ENABLE_HTTPS         | --https         | enable_https   | false    | Enable HTTPS for web server             |
| RedirectHTTP      | REDIRECT_HTTP        | --redirect-http | redirect_http  | true     | Redirect HTTP to HTTPS                  |
| CertDir           | CERT_DIR             | --certdir       | cert_dir       | ./certs  | Directory for certificates              |
| CSRDir            | CSR_DIR              | --csrdir        | csr_dir        | ./csrs   | Directory for CSRs                      |
| LogDir            | LOG_DIR              | --logdir        | log_dir        | ./logs   | Directory for logs                      |
| DatabaseDir       | DB_DIR               | --dbdir         | db_dir         | ./db     | Directory for database                  |

## Extending Configuration

To add new configuration options, simply add new fields to the `Config` struct in `config.go` with appropriate tags:

```go
type Config struct {
    // Existing fields...
    
    // New field
    NewOption string `env:"NEW_OPTION" flag:"new-option" config:"new_option" default:"default_value"`
}
```
