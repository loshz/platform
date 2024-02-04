package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loshz/platform/internal/config"
)

func TestStart(t *testing.T) {
	t.Parallel()

	// Create test Service with required config.
	svc := New("service_test")
	svc.Config().Set(config.KeyServiceStartupTimeout, "1s")

	t.Run("TestStartupError", func(t *testing.T) {
		expected := errors.New("run error")
		runFn := func(context.Context, *Service) error {
			return expected
		}

		err := svc.start(context.Background(), runFn)
		assert.ErrorIs(t, err, expected)
	})

	t.Run("TestNoError", func(t *testing.T) {
		runFn := func(context.Context, *Service) error {
			return nil
		}

		err := svc.start(context.Background(), runFn)
		assert.Nil(t, err)
	})
}

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
