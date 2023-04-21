package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loshz/platform/pkg/config"
)

func TestStart(t *testing.T) {
	// Create test Service with required config.
	svc := New("service_test")
	svc.Config.Set(config.KeyServiceStartupTimeout, "1s")

	t.Run("TestStartupError", func(t *testing.T) {
		expected := errors.New("run error")
		runFn := func(*Service) error {
			return expected
		}

		err := svc.start(runFn)
		assert.ErrorIs(t, err, expected)
	})

	t.Run("TestNoError", func(t *testing.T) {
		runFn := func(*Service) error {
			return nil
		}

		err := svc.start(runFn)
		assert.Nil(t, err)
	})
}

func TestWaitSignal(t *testing.T) {
	svc := New("service_test")

	// Send an error to the service channel.
	go func() {
		svc.errCh <- errors.New("wait error")
	}()

	status := svc.waitSignal()
	assert.Equal(t, status, ExitError)
}

func TestName(t *testing.T) {
	svc := New("service_test")
	assert.Equal(t, svc.Name(), "service_test")
}
