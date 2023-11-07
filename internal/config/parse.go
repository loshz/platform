package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Errors representing parse failures.
var (
	ErrInvalidString      = errors.New("value must be a string")
	ErrInvalidStringSlice = errors.New("value must be a slice of strings")
	ErrInvalidInt         = errors.New("value must be an integer")
	ErrInvalidFloat64     = errors.New("value must be a float64")
	ErrInvalidBool        = errors.New("value must be a boolean")
	ErrInvalidDuration    = errors.New("value must be a duration with a time unit")
	ErrInvalidLogLevel    = errors.New("value must be a log level")
)

// ParseFunc can be used to validate a given configuration value.
// It returns an error if the given value is invalid.
type ParseFunc func(value interface{}) error

// ParseString ensures that a value is a string.
func ParseString(value interface{}) error {
	if value == nil || value == "" {
		return ErrInvalidString
	}

	return nil
}

// ParseStringSlice ensures that a value is comma delimited slice of strings.
// Note: each segment is trimmed.
func ParseStringSlice(value interface{}) error {
	switch value.(type) {
	case string, []string:
		return nil
	}

	return ErrInvalidStringSlice
}

// ParseInt ensures that a value is an int.
func ParseInt(value interface{}) error {
	switch t := value.(type) {
	case int, uint, int32, uint32, int64, uint64:
		return nil
	case string:
		if _, err := strconv.Atoi(t); err == nil {
			return nil
		}
	}

	return ErrInvalidInt
}

// ParseFloat64 ensures that a value is a float64.
func ParseFloat64(value interface{}) error {
	if _, err := strconv.ParseFloat(stringValue(value), 64); err != nil {
		return ErrInvalidFloat64
	}

	return nil
}

// ParseBool ensures that a value is a boolean.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns an error.
func ParseBool(value interface{}) error {
	if _, err := strconv.ParseBool(stringValue(value)); err != nil {
		return ErrInvalidBool
	}

	return nil
}

// ParseDuration ensures that a value is a valid duration.
func ParseDuration(value interface{}) error {
	if _, err := time.ParseDuration(stringValue(value)); err != nil {
		return ErrInvalidDuration
	}

	return nil
}

// ParseLogLevel validates that the value is a valid log level.
// Valid log levels are one of: trace, debug, info, warn, error, fatal
func ParseLogLevel(value interface{}) error {
	switch strings.ToLower(stringValue(value)) {
	case "trace", "debug", "info", "warn", "error", "fatal":
		return nil
	}

	return ErrInvalidLogLevel
}

func stringValue(value interface{}) string {
	switch t := value.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	case nil:
		return ""
	}

	return fmt.Sprintf("%v", value)
}

// stringSliceValues splits a string on commas and trims each segment
func stringSliceValues(s string) []string {
	if strings.TrimSpace(s) == "" {
		return []string{}
	}

	vals := strings.Split(s, ",")
	for i := range vals {
		vals[i] = strings.TrimSpace(vals[i])
	}

	return vals
}
