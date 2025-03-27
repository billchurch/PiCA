package crypto

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// IsYubiKeyPresent checks if a YubiKey is connected to the system
func IsYubiKeyPresent() bool {
	switch runtime.GOOS {
	case "linux":
		// Check for YubiKey devices in Linux
		if _, err := os.Stat("/dev/hidraw0"); err == nil {
			// This is a simple check, but might not be reliable on all systems
			return true
		}
		
		// Check using pcsc_scan
		cmd := exec.Command("pcsc_scan", "-r")
		output, err := cmd.CombinedOutput()
		if err == nil && strings.Contains(string(output), "Yubico") {
			return true
		}
		
		return false
		
	case "darwin":
		// Check for YubiKey devices in macOS
		cmd := exec.Command("system_profiler", "SPUSBDataType")
		output, err := cmd.CombinedOutput()
		if err == nil && strings.Contains(string(output), "Yubico") {
			return true
		}
		
		return false
		
	case "windows":
		// Check for YubiKey devices in Windows
		// This is a simplistic approach and might need refinement
		cmd := exec.Command("wmic", "path", "Win32_PnPEntity", "get", "Caption")
		output, err := cmd.CombinedOutput()
		if err == nil && strings.Contains(string(output), "Yubico") {
			return true
		}
		
		return false
		
	default:
		// Unknown OS, assume no YubiKey
		return false
	}
}

// GetPreferredProviderType returns the preferred provider type based on:
// 1. Environment variable PICA_PROVIDER if set
// 2. Presence of YubiKey if env var is not set
func GetPreferredProviderType() ProviderType {
	// Check environment variable first
	providerEnv := os.Getenv("PICA_PROVIDER")
	switch providerEnv {
	case "software":
		return SoftwareProviderType
	case "yubikey":
		return YubiKeyProviderType
	default:
		// Auto-detect based on YubiKey presence
		if IsYubiKeyPresent() {
			return YubiKeyProviderType
		}
		return SoftwareProviderType
	}
}

// GetProviderNameByType returns a human-readable name for a provider type
func GetProviderNameByType(providerType ProviderType) string {
	switch providerType {
	case SoftwareProviderType:
		return "Software Provider"
	case YubiKeyProviderType:
		return "YubiKey Provider"
	default:
		return "Unknown Provider"
	}
}
