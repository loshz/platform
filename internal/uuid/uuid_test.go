package uuid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	uuid := New("test-service")

	expected := fmt.Sprintf("%s-%s", uuid.Name(), uuid.ID())
	require.Equal(t, expected, uuid.String())
}
