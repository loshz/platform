package config

import (
	"strconv"
	"time"
)

// String attempts to retrieve a config value as a string, or returns a zero value.
func (c *Config) String(key string) string {
	value := c.Get(key)

	return stringValue(value)
}

// StringSlice attempts to retrieve a config value as a slice of strings, or returns an
// empty slice.
func (c *Config) StringSlice(key string) []string {
	value := c.Get(key)

	switch t := value.(type) {
	case []string:
		return t
	case string:
		return stringSliceValues(t)
	}

	return []string{}
}

// Int attempts to retrieve a config value as an int, or returns a zero value.
func (c *Config) Int(key string) int {
	value := c.Get(key)

	switch t := value.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	case string:
		val, err := strconv.Atoi(t)
		if err == nil {
			return val
		}
	}

	return 0
}

// Uint attempts to retrieve a config value as a uint, or returns a zero value.
func (c *Config) Uint(key string) uint {
	value := c.Get(key)

	switch t := value.(type) {
	case uint:
		return uint(t)
	case uint32:
		return uint(t)
	case uint64:
		return uint(t)
	case string:
		val, err := strconv.Atoi(t)
		if err == nil {
			return uint(val)
		}
	}

	return 0
}

// Float64 attempts to retrieve a config value as a float64, or returns a zero value.
func (c *Config) Float64(key string) float64 {
	value := c.Get(key)

	switch t := value.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case string:
		val, err := strconv.ParseFloat(stringValue(value), 64)
		if err == nil {
			return val
		}
	}

	return 0.0
}

// Bool parses a config value as a boolean, or returns false if it cannot converted.
func (c *Config) Bool(key string) bool {
	value := c.Get(key)

	switch t := value.(type) {
	case bool:
		return t
	case string:
		if b, err := strconv.ParseBool(stringValue(value)); err == nil {
			return b
		}
	}

	return false
}

// Duration attempts to parse a config value as a Duration, or returns a zero value.
func (c *Config) Duration(key string) time.Duration {
	value := c.Get(key)

	switch t := value.(type) {
	case time.Duration:
		return t
	case string:
		if dur, err := time.ParseDuration(stringValue(value)); err == nil {
			return dur
		}
	}

	return 0
}
