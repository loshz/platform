package version

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// SemanticVersion is a sequence of three digits used to uniquely identify a software's version
type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

// NewSemanticVersion initializes a SemanticVersion struct from a version string. Strings should be of the form "major.minor.patch", ex: "1.2.3".
func NewSemanticVersion(s string) (*SemanticVersion, error) {
	// older version were of the form `1.0.4-252`, ignore the tag
	dashParts := strings.SplitN(s, "-", 2)

	// might have a leading `v`
	version := strings.TrimLeft(dashParts[0], "v")

	versionParts := strings.Split(version, ".")
	if len(versionParts) != 3 {
		return nil, fmt.Errorf("invalid semantic version: %s", s)
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major part: %s", versionParts[0])
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor part: %s", versionParts[1])
	}

	patch, err := strconv.Atoi(versionParts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch part: %s", versionParts[2])
	}

	return &SemanticVersion{Major: major, Minor: minor, Patch: patch}, nil
}

// MustNewSemanticVersion initializes a SemanticVersion and panics if there's an error
func MustNewSemanticVersion(s string) *SemanticVersion {
	sv, err := NewSemanticVersion(s)
	if err != nil {
		panic(err)
	}
	return sv
}

// Array converts a SemanticVersion struct to a array
func (sv SemanticVersion) Array() [3]int {
	return [3]int{sv.Major, sv.Minor, sv.Patch}
}

// Equal tests two semantic versions for equality
func (sv SemanticVersion) Equal(c SemanticVersion) bool {
	return sv.Major == c.Major &&
		sv.Minor == c.Minor &&
		sv.Patch == c.Patch
}

// GreaterThan compares two semantic versions
// a.GreaterThan(b) -> a > b
func (sv SemanticVersion) GreaterThan(c SemanticVersion) bool {
	base := sv.Array()
	cmp := c.Array()

	for i := range base {
		if base[i] == cmp[i] {
			continue
		}

		if base[i] > cmp[i] {
			return true
		}

		if base[i] < cmp[i] {
			return false
		}
	}

	// equal
	return false
}

// LessThan compares two semantic versions
// a.LessThan(b) -> a < b
func (sv SemanticVersion) LessThan(c SemanticVersion) bool {
	return !sv.GreaterThan(c) && !sv.Equal(c)
}

// GreaterThanEqual compares two semantic versions
// a.GreaterThanEqual(b) -> a >= b
func (sv SemanticVersion) GreaterThanEqual(c SemanticVersion) bool {
	return sv.GreaterThan(c) || sv.Equal(c)
}

// LessThanEqual compares two semantic versions
// a.LessThanEqual(b) -> a <= b
func (sv SemanticVersion) LessThanEqual(c SemanticVersion) bool {
	return sv.LessThan(c) || sv.Equal(c)
}

// InRange checks if a version is within inclusive range of min and max
func (sv SemanticVersion) InRange(min SemanticVersion, max SemanticVersion) bool {
	return sv.GreaterThanEqual(min) && sv.LessThanEqual(max)
}

// String implements the Stringer interface for the SemanticVersion struct
func (sv SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (sv *SemanticVersion) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	v, err := NewSemanticVersion(s)
	if err != nil {
		return err
	}

	*sv = *v
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (sv SemanticVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(sv.String())
}

// SemanticVersions is a list of SemanticVersion
type SemanticVersions []SemanticVersion

// Contains checks if at least one version in the version slice matches value
func (svs SemanticVersions) Contains(value SemanticVersion) bool {
	for _, version := range svs {
		if value.Equal(version) {
			return true
		}
	}
	return false
}

// Max returns the maximum version from the version slice
func (svs SemanticVersions) Max() SemanticVersion {
	var max SemanticVersion
	for _, version := range svs {
		if version.GreaterThan(max) {
			max = version
		}
	}
	return max
}
