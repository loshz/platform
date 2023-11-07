package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Setenv("PLATFORM_SOME_KEY", "value")

	c := New()

	// Test load success.
	err := c.Load("some.key", "default_value", func(interface{}) error {
		return nil
	})
	assert.NoError(t, err)

	// Test default value success.
	err = c.Load("does.not.exist", "default_value", func(interface{}) error {
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, c.Get("does.not.exist"), "default_value")

	// Test validation error.
	err = c.Load("some.key", "default_value", func(interface{}) error {
		return errors.New("validation error")
	})
	assert.EqualError(t, errors.Unwrap(err), "validation error")
}

func TestSetGet(t *testing.T) {
	c := New()
	c.Set("key", "value")

	val := c.Get("key")
	assert.Equal(t, "value", val)
}

func TestNormalizeKey(t *testing.T) {
	key := normalizeKey("some.key")
	assert.Equal(t, key, "PLATFORM_SOME_KEY")

	key = normalizeKey(" some. key ")
	assert.Equal(t, key, "PLATFORM_SOME_KEY")
}
