package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Setenv("TEST_SOME_KEY", "value")

	c := New("test")

	// Test load success.
	err := c.Load("some.key", "", true, func(interface{}) error {
		return nil
	})
	assert.NoError(t, err)

	// Test default value success.
	err = c.Load("does.not.exist", "default_value", true, func(interface{}) error {
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, c.Get("does.not.exist"), "default_value")

	// Test validation error.
	err = c.Load("some.key", "", true, func(interface{}) error {
		return errors.New("validation error")
	})
	assert.EqualError(t, errors.Unwrap(err), "validation error")

	// Test required error.
	err = c.Load("does.not.exist", "", true, func(interface{}) error {
		return nil
	})
	assert.EqualError(t, err, "config value 'TEST_DOES_NOT_EXIST' is required")
}

func TestSetGet(t *testing.T) {
	c := New("test")
	c.Set("key", "value")

	val := c.Get("key")
	assert.Equal(t, "value", val)
}

func TestNormalizeKey(t *testing.T) {
	c := New("test")

	key := c.normalizeKey("some.key")
	assert.Equal(t, key, "TEST_SOME_KEY")

	key = c.normalizeKey(" some. key ")
	assert.Equal(t, key, "TEST_SOME_KEY")
}
