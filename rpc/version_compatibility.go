package rpc

import (
	"fmt"
	"strings"
)

// VersionCompatibilityError represents an error when the RPC provider version is incompatible with the SDK
type VersionCompatibilityError struct {
	ProviderVersion string
	SDKVersion     string
}

// Error returns a formatted error message for version incompatibility
func (e *VersionCompatibilityError) Error() string {
	return fmt.Sprintf("RPC provider version %s is incompatible with SDK version %s. This may cause unexpected behavior.", e.ProviderVersion, e.SDKVersion)
}

// CheckVersionCompatibility checks if the RPC provider version is compatible with the SDK version
//
// Parameters:
//   - providerVersion: The version string returned by the RPC provider
//   - sdkVersion: The version string of the SDK
//
// Returns:
//   - error: A VersionCompatibilityError if versions are incompatible, nil otherwise
func CheckVersionCompatibility(providerVersion, sdkVersion string) error {
	// Parse versions
	providerParts := strings.Split(providerVersion, ".")
	sdkParts := strings.Split(sdkVersion, ".")

	// Validate version format
	if len(providerParts) < 2 || len(sdkParts) < 2 {
		return &VersionCompatibilityError{
			ProviderVersion: providerVersion,
			SDKVersion:     sdkVersion,
		}
	}

	// Compare major and minor versions
	providerMajorMinor := strings.Join(providerParts[:2], ".")
	sdkMajorMinor := strings.Join(sdkParts[:2], ".")

	if providerMajorMinor != sdkMajorMinor {
		return &VersionCompatibilityError{
			ProviderVersion: providerVersion,
			SDKVersion:     sdkVersion,
		}
	}

	return nil
} 