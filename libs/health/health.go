package healthchecks

import (
	"fmt"

	health "github.com/AppsFlyer/go-sundheit"
	"gitlab.com/youwol/platform/libs/go-libs/utils"
	"go.uber.org/zap"
)

// ProbeType indicates if the checker is going to apply for readiness or liveness
type ProbeType int

// ProbeTypeReadiness is the type that applies to the readiness probe
// ProbeTypeLiveness applies to the liveness probe
// ProbeTypeNumber is a facilitator to iterate over types
const (
	ProbeTypeReadiness ProbeType = iota
	ProbeTypeLiveness
	ProbeTypeNumber
)

var healthManagers []health.Health = make([]health.Health, ProbeTypeNumber)

// NewHealthManager creates a new health manager that will continuously check the provided
// checkers configurations
func NewHealthManager(checkers []*health.Config, t ProbeType) (*health.Health, error) {
	if t < ProbeTypeReadiness || t >= ProbeTypeNumber {
		return nil, fmt.Errorf("Incorrect probe type (%d)", t)
	}

	if checkers == nil {
		return nil, fmt.Errorf("Uninitialized checkers list")
	}

	if healthManagers[t] != nil {
		healthManagers[t].DeregisterAll()
	} else {
		healthManagers[t] = health.New()
	}

	log := utils.NewLogger()
	var err error
	for _, c := range checkers {
		err = healthManagers[t].RegisterCheck(c)
		if err != nil {
			log.Error("Cannot register checker", zap.Error(err))
		}
	}

	return &healthManagers[t], err
}

// GetHealth return the health manager associated with the probe type
// func GetHealth(t ProbeType) *health.Health {
// 	return &healthManagers[t]
// }
