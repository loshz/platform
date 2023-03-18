package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testGetterConfig *Config

func TestMain(m *testing.M) {
	// Create test config
	testGetterConfig = New("test")

	// Set string values
	testGetterConfig.Set("string", "some_value")

	// Set string slice values
	testGetterConfig.Set("stringslice1", []string{"a", "b", "c"})
	testGetterConfig.Set("stringslice2", "d,e,f")

	// Set int values
	testGetterConfig.Set("intstring", "1")
	testGetterConfig.Set("int", 1)
	testGetterConfig.Set("uint", uint(1))
	testGetterConfig.Set("int32", int32(1))
	testGetterConfig.Set("uint32", uint32(1))
	testGetterConfig.Set("int64", int64(1))
	testGetterConfig.Set("uint64", uint64(1))

	// Set float64 values
	testGetterConfig.Set("floatstring", "10.00")
	testGetterConfig.Set("float32", float32(10.00))
	testGetterConfig.Set("float64", float64(10.00))

	// Set bool values
	testGetterConfig.Set("bool", true)
	testGetterConfig.Set("booltrue", "true")
	testGetterConfig.Set("boolfalse", "false")

	// Set duration values
	testGetterConfig.Set("duration", time.Second)
	testGetterConfig.Set("durationstring", "1s")

	os.Exit(m.Run())
}

func TestConfigString(t *testing.T) {
	t.Parallel()

	// Assert that a not found key is empty.
	s := testGetterConfig.String("not_found")
	assert.Empty(t, s)

	// Assert that a string key has the expected value.
	s = testGetterConfig.String("string")
	assert.Equal(t, "some_value", s)
}

func TestConfigStringSlice(t *testing.T) {
	t.Parallel()

	// Assert that a not found key is empty.
	s := testGetterConfig.StringSlice("not_found")
	assert.Empty(t, s)

	// Get 1st string slice and compare values.
	s = testGetterConfig.StringSlice("stringslice1")
	assert.Len(t, s, 3)
	assert.Equal(t, []string{"a", "b", "c"}, s)

	// Get 2nd string slice and compare values.
	s = testGetterConfig.StringSlice("stringslice2")
	assert.Len(t, s, 3)
	assert.Equal(t, []string{"d", "e", "f"}, s)

}

func TestConfigInt(t *testing.T) {
	t.Parallel()

	// Assert that a not found key is zero.
	i := testGetterConfig.Int("not_found")
	assert.Empty(t, i)

	// Assert that all int keys return the expected value.
	ints := []string{
		"intstring",
		"int",
		"uint",
		"int32",
		"uint32",
		"int64",
		"uint64",
	}

	for _, in := range ints {
		i := testGetterConfig.Int(in)
		assert.Equal(t, 1, i)
	}
}

func TestConfigFloat64(t *testing.T) {
	t.Parallel()

	// Assert that a not found key is zero.
	f := testGetterConfig.Float64("not_found")
	assert.Empty(t, f)

	// Assert that all float keys return the expected value.
	floats := []string{
		"floatstring",
		"float32",
		"float64",
	}

	for _, fl := range floats {
		f := testGetterConfig.Float64(fl)
		assert.Equal(t, 10.00, f)
	}
}

func TestConfigBool(t *testing.T) {
	t.Parallel()

	// Assert that a not found key is zero.
	b := testGetterConfig.Bool("not_found")
	assert.False(t, b)

	// Assert that all bool keys return the expected value.
	b = testGetterConfig.Bool("bool")
	assert.True(t, b)

	b = testGetterConfig.Bool("booltrue")
	assert.True(t, b)

	b = testGetterConfig.Bool("boolfalse")
	assert.False(t, b)

}

func TestConfigDuration(t *testing.T) {
	t.Parallel()

	// Assert that a not found key is zero.
	d := testGetterConfig.Duration("not_found")
	assert.Empty(t, d)

	// Assert that all duration keys return the expected value.
	d = testGetterConfig.Duration("duration")
	assert.Equal(t, time.Second, d)

	d = testGetterConfig.Duration("durationstring")
	assert.Equal(t, time.Second, d)
}
