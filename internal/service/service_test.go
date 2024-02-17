package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaitSignal(t *testing.T) {
	t.Parallel()

	svc := New("service_test")

	t.Run("TestError", func(t *testing.T) {
		// Send an error to the service channel.
		go func() {
			svc.Error(errors.New("wait error"))
		}()

		status := svc.waitSignal(context.Background())
		assert.Equal(t, status, ExitError)
	})

	t.Run("TestContextDone", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel the context.
		go cancel()

		status := svc.waitSignal(ctx)
		assert.Equal(t, status, ExitOK)
	})
}

func TestName(t *testing.T) {
	t.Parallel()

	// Assert correct name.
	svc := New("service_test")
	assert.Equal(t, svc.Name(), "service_test")

	// Assert name is lowercase.
	svc = New("SERVICE_UPPERCASE")
	assert.Equal(t, svc.Name(), "service_uppercase")
}
