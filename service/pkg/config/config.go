// Package config provides functions to load config from os env vars
// with the ability to perform optional validation on the returned values.
package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// The env var prefix for platform config.
const envVarPrefix = "PCONF_"

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
// The method will return the raw, unparsed value.
func (c *Config) Load(key string, value interface{}, required bool, fns ...ParseFunc) interface{} {
	// read env var from os
	if env := os.Getenv(normalizeKey(key)); env != "" {
		value = env
	}

	// check if value is required
	if value == nil || value == "" {
		if required {
			panic(fmt.Sprintf("value for key '%s' is required", key))
		}

		return nil
	}

	// run validate funtions
	for _, fn := range fns {
		if err := fn(value); err != nil {
			panic(fmt.Sprintf("fatal error validating config value for '%s': %v", key, err))
		}
	}

	c.Set(key, value)

	return value
}

// normalizeKey transforms a config key into a prefixed env var.
// For example: log.level becomes PCONF_LOG_LEVEL
func normalizeKey(key string) string {
	// transform key to env var
	key = strings.Replace(key, ".", "_", -1)
	// add prefix to key
	key = strings.ToUpper(envVarPrefix + key)

	return key
}
