package health_test

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
)

func TestScyllaCheckerOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		resp := make([]healthcheckers.NodeState, 0)
		resp = append(resp, healthcheckers.NodeState{Key: "toto", Value: "UP"})
		r, _ := json.Marshal(resp)
		rw.Write(r)
	}))

	defer server.Close()

	sc := make([]*health.Config, 0)
	sc = append(sc, healthcheckers.NewScyllaChecker(&checks.HTTPCheckConfig{CheckName: "scylla-check-fake", URL: server.URL}))
	h, err := healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeReadiness)
	if err != nil {
		t.Error("Initialization of readiness health manager is not working")
	}

	// Leave some time for the checker to run
	time.Sleep(10 * time.Millisecond)

	_, healthy := (*h).Results()
	if !healthy {
		t.Error("Scylla checker failed")
	}

	// fmt.Printf("Response: %v", res)
}

func TestScyllaCheckerFAIL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		resp := make([]healthcheckers.NodeState, 0)
		resp = append(resp, healthcheckers.NodeState{Key: "toto", Value: "UP"})
		resp = append(resp, healthcheckers.NodeState{Key: "titi", Value: "DOWN"})
		r, _ := json.Marshal(resp)
		rw.Write(r)
	}))

	defer server.Close()

	sc := make([]*health.Config, 0)
	sc = append(sc, healthcheckers.NewScyllaChecker(&checks.HTTPCheckConfig{CheckName: "scylla-check-fake", URL: server.URL}))
	h, err := healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeReadiness)
	if err != nil {
		t.Error("Initialization of readiness health manager is not working")
	}

	// Leave some time for the checker to run
	time.Sleep(10 * time.Millisecond)

	_, healthy := (*h).Results()
	if healthy {
		t.Error("Scylla checker should fail")
	}

	// fmt.Printf("Response: %v", res)
}
