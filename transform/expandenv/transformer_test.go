package expandenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandEnvTransformer_Transform(t *testing.T) {
	t.Parallel()

	envs := map[string]string{
		"TEST_ENV":   "test value",
		"TEST_ENV_1": "test value 1",
	}

	transformer := NewTransformer(WithGetenv(func(s string) string {
		return envs[s]
	}))
	values := map[string]interface{}{
		"test": "${TEST_ENV}",
		"a": []interface{}{
			"test",
			"${TEST_ENV_1}",
		},
	}

	expanded, err := transformer.Transform(values)
	require.NoError(t, err)
	assert.Equal(t, "test value", expanded["test"])

	a, ok := expanded["a"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, "test", a[0])
	assert.Equal(t, "test value 1", a[1])
}

func TestExpandEnvTransformer_Transform_WithDefaultOsGetenv(t *testing.T) {
	t.Parallel()

	transformer := NewTransformer()
	values := map[string]interface{}{
		"test": "${TEST_ENV}",
	}

	expanded, err := transformer.Transform(values)
	require.NoError(t, err)
	assert.Empty(t, expanded["test"])
}
