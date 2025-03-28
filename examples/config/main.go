// Example demonstrating the PiCA configuration system
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/billchurch/PiCA/internal/config"
)

func main() {
	// Load configuration from all sources
	cfg, err := config.Load(os.Args[1:], "")
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Print the loaded configuration
	fmt.Println("PiCA Configuration Example")
	fmt.Println("=========================")
	fmt.Println("")
	
	// Convert to JSON for display
	jsonConfig, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling config to JSON: %v\n", err)
	} else {
		fmt.Println(string(jsonConfig))
	}
	
	// Show configuration sources
	fmt.Println("\nConfiguration Sources Used")
	fmt.Println("-------------------------")
	fmt.Printf("Command-line args: %v\n", len(os.Args) > 1)
	fmt.Printf("Environment vars: PICA_PROVIDER=%s, CA_TYPE=%s\n", 
		os.Getenv("PICA_PROVIDER"), os.Getenv("CA_TYPE"))
	
	// Show actual values from different priority sources
	fmt.Println("\nConfiguration Priority Example")
	fmt.Println("-----------------------------")
	fmt.Printf("CA Type: %s (Priority: Command-line > Environment > Config file > Default)\n", cfg.CAType)
	fmt.Printf("Provider: %s (Priority: Command-line > Environment > Config file > Default)\n", cfg.ProviderType)
	fmt.Printf("Web Port: %d (Priority: Command-line > Environment > Config file > Default)\n", cfg.WebPort)
	
	// Show usage instructions
	fmt.Println("\nTry running with environment variables:")
	fmt.Println("  export CA_TYPE=root")
	fmt.Println("  export PICA_PROVIDER=software")
	fmt.Println("  go run examples/config/main.go")
	fmt.Println("")
	fmt.Println("Or with command-line arguments:")
	fmt.Println("  go run examples/config/main.go --ca-type=root --provider=yubikey --port=8443")
	fmt.Println("")
	fmt.Println("Or with a configuration file:")
	fmt.Println("  go run examples/config/main.go --config=./configs/examples/root-ca.json")
}
