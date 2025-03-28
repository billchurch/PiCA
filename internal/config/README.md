# PiCA Configuration Package

This package provides a unified configuration system for PiCA with support for:

1. Command-line flags (highest priority)
2. Environment variables 
3. Configuration files (JSON/TOML)
4. Default values (lowest priority)

The package also handles validation of configuration values to ensure required fields are present and values are properly formatted.

## Basic Usage

```go
import (
    "github.com/billchurch/PiCA/internal/config"
    "os"
)

func main() {
    // Load configuration from all sources
    cfg, err := config.Load(os.Args[1:], "")
    if err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }
    
    // Use configuration values
    fmt.Printf("CA Type: %s\n", cfg.CAType)
    fmt.Printf("Web Port: %d\n", cfg.WebPort)
}
```

## Adding New Options

To add a new configuration option:

1. Add a field to the `Config` struct in `config.go`
2. Include field tags for environment variable, command-line flag, config file, and default value
3. The option will be automatically handled by the loading system

Example:

```go
// New field in Config struct
MaxConnections int `env:"MAX_CONNECTIONS" flag:"max-connections" config:"max_connections" default:"100"`
```

## Validation

The configuration system validates settings after loading from all sources. The `Validate()` method checks:

1. Required fields based on application type (e.g., CA config and cert files required for web server)
2. YubiKey slot format (must be valid hex)
3. HTTPS settings (TLS certificate and key files required when HTTPS is enabled)

Validation errors are returned from the `Load()` function, preventing the application from starting with invalid configuration.

## Testing

Run tests with:

```
go test -v ./internal/config
```

For more details on the configuration system, see the documentation in [../docs/configuration.md](../docs/configuration.md).
