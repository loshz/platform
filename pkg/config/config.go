// Package config provides functions to load config from os env vars
// with the ability to perform optional validation on the returned values.
package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Config stores config key/values and provides methods
// for concurrent read/write access.
type Config struct {
	values map[string]interface{}
	mtx    sync.RWMutex
}

// New creates a new config with an initialized in memory store.
func New() *Config {
	c := &Config{
		values: make(map[string]interface{}),
	}

	return c
}

// Get finds a config value by key.
func (c *Config) Get(key string) interface{} {
	c.mtx.RLock()
	val := c.values[key]
	c.mtx.RUnlock()

	return val
}

// Set sets a config value by key.
func (c *Config) Set(key string, value interface{}) {
	c.mtx.Lock()
	c.values[key] = value
	c.mtx.Unlock()
}

// Load attempts to read config values from env vars, setting a default value if not found.
// If supplied, all parse funcs will be ran against the value and panic on failure.
func (c *Config) Load(key string, value interface{}, fns ...ParseFunc) error {
	normKey := normalizeKey(key)

	// Read valud from env vars.
	if env := os.Getenv(normKey); env != "" {
		value = env
	}

	// Check for default value.
	if value == nil || value == "" {
		return fmt.Errorf("error: required config value '%s' not set", normKey)
	}

	// Run validate funtions.
	for _, fn := range fns {
		if err := fn(value); err != nil {
			return fmt.Errorf("error validating config value for '%s': %w", normKey, err)
		}
	}

	c.Set(key, value)

	return nil
}

// MustLoad is functionally equivalent to Load, but panics on error.
func (c *Config) MustLoad(key string, value interface{}, fns ...ParseFunc) {
	if err := c.Load(key, value, fns...); err != nil {
		panic(err)
	}
}

// normalizeKey transforms a config key into a prefixed env var.
// For example: log.level becomes PLATFORM_LOG_LEVEL
func normalizeKey(key string) string {
	// Remove whiitespace.
	key = strings.ReplaceAll(key, " ", "")
	// Transform key to env var.
	key = strings.Replace(key, ".", "_", -1)
	// Add prefix to key.
	key = strings.ToUpper(fmt.Sprintf("PLATFORM_%s", key))

	return key
}
