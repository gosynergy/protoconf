package protoconf_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gosynergy/protoconf/transform/expandenv"
)

func TestExpandEnvTransformer_Transform(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		transformer := expandenv.New(expandenv.WithGetenv(func(s string) string {
			return "test value"
		}))
		values := map[string]interface{}{
			"test": "${TEST_ENV}",
		}

		expanded, err := transformer.Transform(values)
		require.NoError(t, err)
		assert.Equal(t, "test value", expanded["test"])
	})
}
