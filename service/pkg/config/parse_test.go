package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	t.Parallel()

	// Assert invalid string returns an error.
	err := ParseString(nil)
	assert.ErrorIs(t, err, ErrInvalidString)

	// Assert valid string does not return an error.
	err = ParseString("some_value")
	assert.ErrorIs(t, err, nil)
}

func TestParseStringSlice(t *testing.T) {
	t.Parallel()

	// Assert invalid string slice returns an error.
	err := ParseStringSlice(1)
	assert.ErrorIs(t, err, ErrInvalidStringSlice)

	// Assert valid string slices do not return an error.
	err = ParseStringSlice("slice,of,strings")
	assert.ErrorIs(t, err, nil)

	err = ParseStringSlice([]string{"slice", "of", "strings"})
	assert.ErrorIs(t, err, nil)
}

func TestParseInt(t *testing.T) {
	t.Parallel()

	// Assert invalid int returns an error.
	err := ParseInt("invalid")
	assert.ErrorIs(t, err, ErrInvalidInt)

	// Assert that all int keys return no error.
	ints := []interface{}{
		"1",
		int(1),
		uint(1),
		int32(1),
		uint32(1),
		int64(1),
		uint64(1),
	}

	for _, i := range ints {
		err := ParseInt(i)
		assert.ErrorIs(t, err, nil)
	}
}

func TestParseFloat64(t *testing.T) {
	t.Parallel()

	// Assert invalid float returns an error.
	err := ParseFloat64("invalid")
	assert.ErrorIs(t, err, ErrInvalidFloat64)

	// Assert valid float does not return an error.
	err = ParseFloat64("11.99")
	assert.ErrorIs(t, err, nil)
}

func TestParseBool(t *testing.T) {
	t.Parallel()

	// Assert invalid bool returns an error.
	err := ParseBool("invalid")
	assert.ErrorIs(t, err, ErrInvalidBool)

	// Assert that all bool keys return no error.
	bools := []interface{}{
		1,
		"1",
		"t",
		"T",
		"TRUE",
		"true",
		"True",
		true,
		0,
		"0",
		"f",
		"F",
		"FALSE",
		"false",
		"False",
		false,
	}

	for _, b := range bools {
		err := ParseBool(b)
		assert.ErrorIs(t, err, nil)
	}
}

func TestParseDuration(t *testing.T) {
	t.Parallel()

	// Assert invalid duration returns an error.
	err := ParseDuration("invalid")
	assert.ErrorIs(t, err, ErrInvalidDuration)

	// Assert that all duration keys return no error.
	durations := []interface{}{
		"1h",
		"10s",
		"30m",
		"500ms",
		"2h45m",
	}

	for _, dur := range durations {
		err := ParseDuration(dur)
		assert.ErrorIs(t, err, nil)
	}
}

func TestParseLogLevel(t *testing.T) {
	t.Parallel()

	// Assert invalid duration returns an error.
	err := ParseLogLevel("invalid")
	assert.ErrorIs(t, err, ErrInvalidLogLevel)

	// Assert that all log level keys return no error.
	levels := []interface{}{
		"trace",
		"debug",
		"info",
		"warn",
		"error",
		"fatal",
	}

	for _, level := range levels {
		err := ParseLogLevel(level)
		assert.ErrorIs(t, err, nil)
	}
}
