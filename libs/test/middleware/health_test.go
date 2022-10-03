package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	healthchecks "gitlab.com/youwol/platform/libs/go-libs/health"
	healthcheckers "gitlab.com/youwol/platform/libs/go-libs/health/checkers"
	middleware "gitlab.com/youwol/platform/libs/go-libs/middleware"
)

func TestReadinessMiddleware(t *testing.T) {
	// First server to provide the checker endpoint
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		resp := make([]healthcheckers.NodeState, 0)
		resp = append(resp, healthcheckers.NodeState{Key: "toto", Value: "UP"})
		r, _ := json.Marshal(resp)
		rw.Write(r)
	}))

	defer server.Close()

	sc := make([]*health.Config, 0)
	sc = append(sc,
		healthcheckers.NewScyllaChecker(
			&checks.HTTPCheckConfig{
				CheckName: "scylla-check-fake",
				URL:       server.URL,
			}))

	h := middleware.NewHealthMiddleware(nil, sc, healthchecks.ProbeTypeReadiness)

	// Leave some time for the checker to run
	time.Sleep(10 * time.Millisecond)
	req := httptest.NewRequest("GET", "/ready", nil)
	recorder := httptest.NewRecorder() // Second server to provide the /ready endpoint
	h.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("Readiness probe returned %d", recorder.Code)
	}
}

func TestLivenessMiddleware(t *testing.T) {
	// First server to provide the checker endpoint
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		resp := make([]healthcheckers.NodeState, 0)
		resp = append(resp, healthcheckers.NodeState{Key: "toto", Value: "UP"})
		r, _ := json.Marshal(resp)
		rw.Write(r)
	}))

	defer server.Close()

	sc := make([]*health.Config, 0)
	sc = append(sc,
		healthcheckers.NewScyllaChecker(
			&checks.HTTPCheckConfig{
				CheckName: "scylla-check-fake",
				URL:       server.URL,
			}))

	h := middleware.NewHealthMiddleware(nil, sc, healthchecks.ProbeTypeLiveness)

	// Leave some time for the checker to run
	time.Sleep(10 * time.Millisecond)

	req := httptest.NewRequest("GET", "/alive", nil)
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("Liveness probe returned %d", recorder.Code)
	}
}
