package config

import (
	"strconv"
	"time"
)

// String retrieves a config value as a string, or returns zero value
// if it cannot be converted.
func (c *Config) String(key string) string {
	value := c.Get(key)

	return stringValue(value)
}

// StringSlice retrieves a config value as a slice of string, or returns an
// empty slice (not nil) if it cannot converted.
func (c *Config) StringSlice(key string) []string {
	value := c.Get(key)

	switch t := value.(type) {
	case []string:
		return t
	case string:
		return stringSliceValues(t)
	}

	return nil
}

// Int retrieves a config value as an int, or returns zero value
// if it cannot be converted to an int.
func (c *Config) Int(key string) int {
	value := c.Get(key)

	switch t := value.(type) {
	case int:
		return t
	case uint:
		return int(t)
	case int32:
		return int(t)
	case uint32:
		return int(t)
	case int64:
		return int(t)
	case uint64:
		return int(t)
	case string:
		val, err := strconv.Atoi(t)
		if err == nil {
			return val
		}
	}

	return 0
}

// Float64 retrieves a config value as a float64, or returns zero value
// if it cannot be converted.
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

// Duration parses a config value as a Duration, or returns zero value if it
// cannot be converted to a Duration.
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
