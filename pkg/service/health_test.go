package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockedService = "mocked-service"

type mockHealthCheckedDependency struct {
	name string
	err  bool
}

func (m mockHealthCheckedDependency) Name() string {
	return m.name
}

func (m mockHealthCheckedDependency) Check() error {
	if m.err {
		return fmt.Errorf("%s error", mockedService)
	}
	return nil
}

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	t.Run("TestLivenessOK", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := healthHandler(mockedService, nil)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "{\"service\":\"mocked-service\",\"status\":\"OK\"}\n", rr.Body.String())
	})

	t.Run("TestReadinessOK", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodGet, "/health?status=ready", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := healthHandler(mockedService, nil)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "{\"service\":\"mocked-service\",\"status\":\"OK\"}\n", rr.Body.String())
	})

	t.Run("TestReadinessFailing", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodGet, "/health?status=ready", nil)
		require.NoError(t, err)

		// create mock services with a forced failure on one
		services := []HealthCheckedDependency{
			&mockHealthCheckedDependency{name: "service-1", err: false},
			&mockHealthCheckedDependency{name: "service-2", err: true},
			&mockHealthCheckedDependency{name: "service-3", err: false},
		}

		rr := httptest.NewRecorder()
		handler := healthHandler(mockedService, services)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

		b, err := io.ReadAll(rr.Body)
		require.NoError(t, err)

		res := &healthResponse{}
		err = json.Unmarshal(b, res)
		require.NoError(t, err)

		assert.Equal(t, res.Status, HealthStatusFailing)

		for _, service := range res.Dependencies {
			if service.Service == "service-1" || service.Service == "service-3" {
				assert.Equal(t, service.Status, HealthStatusOK)
			}
			if service.Service == "service-2" {
				assert.Equal(t, service.Status, HealthStatusFailing)
			}
		}
	})
}
