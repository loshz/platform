package uuid

import (
	"fmt"
	"strings"

	guuid "github.com/google/uuid"
)

// UUID represents an unique service identifier including a name prefix.
// E.g., service-xxxx-xxxx
type UUID struct {
	name string
	id   guuid.UUID
}

func New(name string) UUID {
	if name == "" {
		panic("must provide a name when generating a UUID")
	}

	return UUID{strings.ToLower(name), guuid.New()}
}

func (u UUID) ID() string     { return u.id.String() }
func (u UUID) Name() string   { return u.name }
func (u UUID) String() string { return fmt.Sprintf("%s-%s", u.name, u.id) }
