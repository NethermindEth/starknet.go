package paymaster

import (
	"context"
	"strconv"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Test the 'paymaster_isAvailable' method
//
//nolint:tparallel // Each subtest runs in different environments
func TestIsAvailable(t *testing.T) {
	t.Parallel()
	t.Run("integration", func(t *testing.T) {
		tests.RunTestOn(t, tests.IntegrationEnv)

		pm, spy := SetupPaymaster(t)
		available, err := pm.IsAvailable(context.Background())
		require.NoError(t, err)

		assert.Equal(t, string(spy.LastResponse()), strconv.FormatBool(available))
		assert.True(t, available)
	})

	t.Run("mock", func(t *testing.T) {
		tests.RunTestOn(t, tests.MockEnv)

		pm := SetupMockPaymaster(t)
		pm.c.EXPECT().
			CallContextWithSliceArgs(context.Background(), gomock.AssignableToTypeOf(new(bool)), "paymaster_isAvailable").
			SetArg(1, true).
			Return(nil)
		available, err := pm.IsAvailable(context.Background())
		assert.NoError(t, err)
		assert.True(t, available)
	})
}
