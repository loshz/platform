package version

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newSemVer(s string) *SemanticVersion {
	sv, _ := NewSemanticVersion(s)
	return sv
}

func TestSemanticVersionComparisons(t *testing.T) {
	type testCase struct {
		base           *SemanticVersion
		operator       string
		cmp            *SemanticVersion
		expectedOutput bool
	}

	// might have gotten carried away here, just decided to brute force it:
	// major < major, minor < minor, patch < patch
	// major < major, minor < minor, patch = patch
	// major < major, minor < minor, patch > patch
	// major < major, minor = minor, patch < patch
	// major < major, minor = minor, patch = patch
	// ...
	cases := []testCase{
		// less than
		{newSemVer("0.0.0"), "<", newSemVer("1.1.1"), true},
		{newSemVer("0.0.1"), "<", newSemVer("1.1.1"), true},
		{newSemVer("0.0.1"), "<", newSemVer("1.1.0"), true},
		{newSemVer("0.1.0"), "<", newSemVer("1.1.1"), true},
		{newSemVer("0.1.1"), "<", newSemVer("1.1.1"), true},
		{newSemVer("0.1.1"), "<", newSemVer("1.1.0"), true},
		{newSemVer("0.1.0"), "<", newSemVer("1.0.1"), true},
		{newSemVer("0.1.1"), "<", newSemVer("1.0.1"), true},
		{newSemVer("0.1.1"), "<", newSemVer("1.0.0"), true},
		{newSemVer("1.0.0"), "<", newSemVer("1.1.1"), true},
		{newSemVer("1.0.1"), "<", newSemVer("1.1.1"), true},
		{newSemVer("1.0.1"), "<", newSemVer("1.1.0"), true},
		{newSemVer("1.1.0"), "<", newSemVer("1.1.1"), true},
		{newSemVer("1.1.1"), "<", newSemVer("1.1.1"), false},
		{newSemVer("1.1.1"), "<", newSemVer("1.1.0"), false},
		{newSemVer("1.1.0"), "<", newSemVer("1.0.1"), false},
		{newSemVer("1.1.1"), "<", newSemVer("1.0.1"), false},
		{newSemVer("1.1.1"), "<", newSemVer("1.0.0"), false},
		{newSemVer("1.0.0"), "<", newSemVer("0.1.1"), false},
		{newSemVer("1.0.1"), "<", newSemVer("0.1.1"), false},
		{newSemVer("1.0.1"), "<", newSemVer("0.1.0"), false},
		{newSemVer("1.1.0"), "<", newSemVer("0.1.1"), false},
		{newSemVer("1.1.1"), "<", newSemVer("0.1.1"), false},
		{newSemVer("1.1.1"), "<", newSemVer("0.1.0"), false},
		{newSemVer("1.1.0"), "<", newSemVer("0.0.1"), false},
		{newSemVer("1.1.1"), "<", newSemVer("0.0.1"), false},
		{newSemVer("1.1.1"), "<", newSemVer("0.0.0"), false},

		// equal
		{newSemVer("0.0.0"), "=", newSemVer("1.1.1"), false},
		{newSemVer("0.0.1"), "=", newSemVer("1.1.1"), false},
		{newSemVer("0.0.1"), "=", newSemVer("1.1.0"), false},
		{newSemVer("0.1.0"), "=", newSemVer("1.1.1"), false},
		{newSemVer("0.1.1"), "=", newSemVer("1.1.1"), false},
		{newSemVer("0.1.1"), "=", newSemVer("1.1.0"), false},
		{newSemVer("0.1.0"), "=", newSemVer("1.0.1"), false},
		{newSemVer("0.1.1"), "=", newSemVer("1.0.1"), false},
		{newSemVer("0.1.1"), "=", newSemVer("1.0.0"), false},
		{newSemVer("1.0.0"), "=", newSemVer("1.1.1"), false},
		{newSemVer("1.0.1"), "=", newSemVer("1.1.1"), false},
		{newSemVer("1.0.1"), "=", newSemVer("1.1.0"), false},
		{newSemVer("1.1.0"), "=", newSemVer("1.1.1"), false},
		{newSemVer("1.1.1"), "=", newSemVer("1.1.1"), true},
		{newSemVer("1.1.1"), "=", newSemVer("1.1.0"), false},
		{newSemVer("1.1.0"), "=", newSemVer("1.0.1"), false},
		{newSemVer("1.1.1"), "=", newSemVer("1.0.1"), false},
		{newSemVer("1.1.1"), "=", newSemVer("1.0.0"), false},
		{newSemVer("1.0.0"), "=", newSemVer("0.1.1"), false},
		{newSemVer("1.0.1"), "=", newSemVer("0.1.1"), false},
		{newSemVer("1.0.1"), "=", newSemVer("0.1.0"), false},
		{newSemVer("1.1.0"), "=", newSemVer("0.1.1"), false},
		{newSemVer("1.1.1"), "=", newSemVer("0.1.1"), false},
		{newSemVer("1.1.1"), "=", newSemVer("0.1.0"), false},
		{newSemVer("1.1.0"), "=", newSemVer("0.0.1"), false},
		{newSemVer("1.1.1"), "=", newSemVer("0.0.1"), false},
		{newSemVer("1.1.1"), "=", newSemVer("0.0.0"), false},

		// greater than
		{newSemVer("0.0.0"), ">", newSemVer("1.1.1"), false},
		{newSemVer("0.0.1"), ">", newSemVer("1.1.1"), false},
		{newSemVer("0.0.1"), ">", newSemVer("1.1.0"), false},
		{newSemVer("0.1.0"), ">", newSemVer("1.1.1"), false},
		{newSemVer("0.1.1"), ">", newSemVer("1.1.1"), false},
		{newSemVer("0.1.1"), ">", newSemVer("1.1.0"), false},
		{newSemVer("0.1.0"), ">", newSemVer("1.0.1"), false},
		{newSemVer("0.1.1"), ">", newSemVer("1.0.1"), false},
		{newSemVer("0.1.1"), ">", newSemVer("1.0.0"), false},
		{newSemVer("1.0.0"), ">", newSemVer("1.1.1"), false},
		{newSemVer("1.0.1"), ">", newSemVer("1.1.1"), false},
		{newSemVer("1.0.1"), ">", newSemVer("1.1.0"), false},
		{newSemVer("1.1.0"), ">", newSemVer("1.1.1"), false},
		{newSemVer("1.1.1"), ">", newSemVer("1.1.1"), false},
		{newSemVer("1.1.1"), ">", newSemVer("1.1.0"), true},
		{newSemVer("1.1.0"), ">", newSemVer("1.0.1"), true},
		{newSemVer("1.1.1"), ">", newSemVer("1.0.1"), true},
		{newSemVer("1.1.1"), ">", newSemVer("1.0.0"), true},
		{newSemVer("1.0.0"), ">", newSemVer("0.1.1"), true},
		{newSemVer("1.0.1"), ">", newSemVer("0.1.1"), true},
		{newSemVer("1.0.1"), ">", newSemVer("0.1.0"), true},
		{newSemVer("1.1.0"), ">", newSemVer("0.1.1"), true},
		{newSemVer("1.1.1"), ">", newSemVer("0.1.1"), true},
		{newSemVer("1.1.1"), ">", newSemVer("0.1.0"), true},
		{newSemVer("1.1.0"), ">", newSemVer("0.0.1"), true},
		{newSemVer("1.1.1"), ">", newSemVer("0.0.1"), true},
		{newSemVer("1.1.1"), ">", newSemVer("0.0.0"), true},
	}

	for _, tc := range cases {
		msg := fmt.Sprintf("%s %s %s -> %t", tc.base, tc.operator, tc.cmp, tc.expectedOutput)
		t.Run(msg, func(t *testing.T) {
			switch tc.operator {
			case ">":
				assert.Equal(t, tc.expectedOutput, tc.base.GreaterThan(*tc.cmp), msg)
			case "<":
				assert.Equal(t, tc.expectedOutput, tc.base.LessThan(*tc.cmp), msg)
			case "=":
				assert.Equal(t, tc.expectedOutput, tc.base.Equal(*tc.cmp), msg)
			default:
				t.Fatalf("invalid operator: %s", tc.operator)
			}
		})
	}
}

func TestSemanticVersion_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Major int
		Minor int
		Patch int
	}
	tests := []struct {
		name    string
		fields  fields
		data    []byte
		wantErr bool
	}{
		{"standard", fields{1, 2, 3}, []byte(`"1.2.3"`), false},
		{"with v", fields{1, 2, 3}, []byte(`"v1.2.3"`), false},
		{"with build number", fields{1, 2, 3}, []byte(`"1.2.3-123"`), false},
		{"empty", fields{0, 0, 0}, []byte(`""`), true},
		{"too many pieces", fields{0, 0, 0}, []byte(`"1.2.3.4"`), true},
		{"invalid", fields{0, 0, 0}, []byte(`"1.2.3a"`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sv SemanticVersion
			if err := json.Unmarshal(tt.data, &sv); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.fields.Major, sv.Major)
			assert.Equal(t, tt.fields.Minor, sv.Minor)
			assert.Equal(t, tt.fields.Patch, sv.Patch)
		})
	}
}

func TestSemanticVersion_MarshalJSON(t *testing.T) {
	type fields struct {
		Major int
		Minor int
		Patch int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"regular", fields{1, 2, 3}, []byte(`"1.2.3"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv := SemanticVersion{
				Major: tt.fields.Major,
				Minor: tt.fields.Minor,
				Patch: tt.fields.Patch,
			}
			got, err := sv.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
