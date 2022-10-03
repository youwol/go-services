package healthchecks

import (
	"time"

	health "github.com/AppsFlyer/go-sundheit"
)

type okChecker struct {
	message string
}

// Name returns the checker name
func (check *okChecker) Name() string {
	return "ok-checker"
}

// Execute runs the checker
func (check *okChecker) Execute() (details interface{}, err error) {
	return check.message, err
}

// NewOkChecker creates a health check for any HTTP endpoint
func NewOkChecker(msg string) *health.Config {
	// Instantiate a ok checker config
	return &health.Config{
		Check:            &okChecker{message: msg},
		ExecutionPeriod:  time.Duration(2) * time.Second,
		InitialDelay:     0,
		InitiallyPassing: false,
	}
}
