package health_test

import (
	"testing"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	healthchecks "gitlab.com/youwol/platform/libs/go-libs/health"
	healthcheckers "gitlab.com/youwol/platform/libs/go-libs/health/checkers"
)

func TestHealthManagerBadArguments(t *testing.T) {
	_, err := healthchecks.NewHealthManager(nil, healthchecks.ProbeTypeNumber)
	if err == nil {
		t.Error("Bad arguments should make the function fail")
	}
}

func TestHealthManagerEmptyCheckers(t *testing.T) {
	_, err := healthchecks.NewHealthManager(nil, healthchecks.ProbeTypeNumber)
	if err == nil {
		t.Error("Empty checkers should make the function fail")
	}
}

func TestHealthManagerBadProbeType(t *testing.T) {

	sc := make([]*health.Config, 0)
	sc = append(sc, healthcheckers.NewScyllaChecker(&checks.HTTPCheckConfig{}))
	_, err := healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeNumber)
	if err == nil {
		t.Error("Bad probe type should make the function fail")
	}
}

func TestHealthManagerOk(t *testing.T) {

	sc := make([]*health.Config, 0)
	sc = append(sc, healthcheckers.NewScyllaChecker(&checks.HTTPCheckConfig{CheckName: "scylla-check-fake", URL: "https://duckduckgo.com"}))
	_, err := healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeReadiness)
	if err != nil {
		t.Error("Initialization of health manager is not working")
	}
}

func TestHealthManagerInitializedTwice(t *testing.T) {

	sc := make([]*health.Config, 0)
	sc = append(sc, healthcheckers.NewScyllaChecker(&checks.HTTPCheckConfig{CheckName: "scylla-check-fake", URL: "https://duckduckgo.com"}))
	_, err := healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeReadiness)
	if err != nil {
		t.Error("First initialization of health manager is not working")
	}

	// We must support a second initialization that will overwrite the first one
	_, err = healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeReadiness)
	if err != nil {
		t.Error("Second initialization of health manager is not working")
	}
}

func TestHealthManagerInitializedReadinessAndLiveness(t *testing.T) {

	sc := make([]*health.Config, 0)
	sc = append(sc, healthcheckers.NewScyllaChecker(&checks.HTTPCheckConfig{CheckName: "scylla-check-fake", URL: "https://duckduckgo.com"}))
	_, err := healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeReadiness)
	if err != nil {
		t.Error("Initialization of readiness health manager is not working")
	}

	_, err = healthchecks.NewHealthManager(sc, healthchecks.ProbeTypeLiveness)
	if err != nil {
		t.Error("Initialization of liveness health manager is not working")
	}
}
