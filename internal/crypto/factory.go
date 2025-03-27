package crypto

import (
	"fmt"
)

// DefaultProviderType determines the default provider type to use
func DefaultProviderType() ProviderType {
	return GetPreferredProviderType()
}

// CreateDefaultProvider creates a provider of the default type
func CreateDefaultProvider() (Provider, error) {
	providerType := DefaultProviderType()
	
	// Default options for each provider type
	var opts map[string]interface{}
	
	switch providerType {
	case SoftwareProviderType:
		// Use default options for software provider
		opts = map[string]interface{}{
			"name": "Default Software Provider",
		}
	case YubiKeyProviderType:
		// Use default options for YubiKey provider
		opts = map[string]interface{}{
			"name": "Default YubiKey Provider",
		}
	default:
		return nil, fmt.Errorf("unsupported default provider type: %v", providerType)
	}
	
	provider, err := NewProvider(providerType, opts)
	if err != nil {
		return nil, err
	}
	
	// Connect to the provider
	if err := provider.Connect(); err != nil {
		// If the preferred provider fails to connect (e.g., YubiKey not available),
		// try falling back to the software provider
		if providerType == YubiKeyProviderType {
			fmt.Println("Warning: Failed to connect to YubiKey provider, falling back to software provider")
			return CreateProviderFromConfig(map[string]interface{}{
				"type": "software",
				"name": "Fallback Software Provider",
			})
		}
		return nil, fmt.Errorf("failed to connect to provider: %w", err)
	}
	
	return provider, nil
}

// CreateProviderFromConfig creates a provider based on a configuration
func CreateProviderFromConfig(config map[string]interface{}) (Provider, error) {
	// Extract provider type from config
	providerTypeStr, ok := config["type"].(string)
	if !ok {
		// Use default provider if not specified
		return CreateDefaultProvider()
	}
	
	var providerType ProviderType
	switch providerTypeStr {
	case "software":
		providerType = SoftwareProviderType
	case "yubikey":
		providerType = YubiKeyProviderType
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerTypeStr)
	}
	
	// Create provider with config options
	provider, err := NewProvider(providerType, config)
	if err != nil {
		return nil, err
	}
	
	// Connect to the provider
	if err := provider.Connect(); err != nil {
		// If this is a YubiKey provider and it failed to connect,
		// we could try falling back to software provider
		if providerType == YubiKeyProviderType && config["fallback"] != "false" {
			fmt.Println("Warning: Failed to connect to YubiKey provider, falling back to software provider")
			return CreateProviderFromConfig(map[string]interface{}{
				"type": "software",
				"name": "Fallback Software Provider",
			})
		}
		return nil, fmt.Errorf("failed to connect to provider: %w", err)
	}
	
	return provider, nil
}
