package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("TestSetGet", func(t *testing.T) {
		t.Parallel()

		c := New()
		c.Set("key", "value")

		val := c.Get("key")
		assert.Equal(t, "value", val)
	})
}

func TestNormalizeKey(t *testing.T) {
	t.Parallel()

	expected := "PCONF_LOG_LEVEL"
	if key := normalizeKey(KeyLogLevel); key != expected {
		t.Errorf("expected key: '%s', got: '%s'", expected, key)
	}
}
