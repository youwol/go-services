package healthchecks

import (
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
)

// NewHTTPChecker creates a health check for any HTTP endpoint
func NewHTTPChecker(cfg *checks.HTTPCheckConfig) *health.Config {

	// Instantiate a generic HTTP checker
	ht, _ := checks.NewHTTPCheck(*cfg)
	return &health.Config{
		Check:            ht,
		ExecutionPeriod:  time.Duration(2) * time.Second,
		InitialDelay:     0,
		InitiallyPassing: false,
	}
}
