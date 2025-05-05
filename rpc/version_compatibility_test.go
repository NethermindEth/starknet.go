package rpc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCheckVersionCompatibility tests the version compatibility check functionality
func TestCheckVersionCompatibility(t *testing.T) {
	tests := []struct {
		name            string
		providerVersion string
		sdkVersion     string
		wantErr        bool
		errMsg         string
	}{
		{
			name:            "matching versions",
			providerVersion: "0.8.0",
			sdkVersion:     "0.8.0",
			wantErr:        false,
		},
		{
			name:            "matching major.minor different patch",
			providerVersion: "0.8.1",
			sdkVersion:     "0.8.0",
			wantErr:        false,
		},
		{
			name:            "different major version",
			providerVersion: "1.8.0",
			sdkVersion:     "0.8.0",
			wantErr:        true,
			errMsg:         "RPC provider version 1.8.0 is incompatible with SDK version 0.8.0",
		},
		{
			name:            "different minor version",
			providerVersion: "0.7.0",
			sdkVersion:     "0.8.0",
			wantErr:        true,
			errMsg:         "RPC provider version 0.7.0 is incompatible with SDK version 0.8.0",
		},
		{
			name:            "invalid provider version format",
			providerVersion: "invalid",
			sdkVersion:     "0.8.0",
			wantErr:        true,
			errMsg:         "RPC provider version invalid is incompatible with SDK version 0.8.0",
		},
		{
			name:            "invalid SDK version format",
			providerVersion: "0.8.0",
			sdkVersion:     "invalid",
			wantErr:        true,
			errMsg:         "RPC provider version 0.8.0 is incompatible with SDK version invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckVersionCompatibility(tt.providerVersion, tt.sdkVersion)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
				// Verify it's the correct error type
				_, ok := err.(*VersionCompatibilityError)
				require.True(t, ok, "error should be of type VersionCompatibilityError")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVersionCompatibilityError(t *testing.T) {
	err := &VersionCompatibilityError{
		ProviderVersion: "0.7.0",
		SDKVersion:     "0.8.0",
	}

	expectedMsg := "RPC provider version 0.7.0 is incompatible with SDK version 0.8.0. This may cause unexpected behavior."
	require.Equal(t, expectedMsg, err.Error())
} 