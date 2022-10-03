package healthchecks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/pkg/errors"
	"gitlab.com/youwol/platform/libs/go-libs/utils"
	"go.uber.org/zap"
)

type scyllaChecker struct {
	config         *checks.HTTPCheckConfig
	successDetails string
}

// NodeState is the expected format for scylla nodes statuss
type NodeState struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (check *scyllaChecker) Name() string {
	return check.config.CheckName
}

func (check *scyllaChecker) Execute() (details interface{}, err error) {
	details = check.config.URL

	// prepare the request
	req, err := http.NewRequest(check.config.Method, check.config.URL, nil)
	l := utils.NewLogger()
	if err != nil {
		l.Error("Unable to create request for scylla checker", zap.Error(err))
		return nil, err
	}
	for _, opt := range check.config.Options {
		opt(req)
	}

	// execute request
	resp, err := check.config.Client.Do(req)
	if err != nil {
		l.Error("Unable to execute request for scylla checker", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Checker returned HTTP status %d", resp.StatusCode)
		l.Error("Scylla checker error", zap.Error(err))
		return details, err
	}

	// Check that all nodes are UP
	body, err := ioutil.ReadAll(resp.Body)
	s := string(body)
	l.Debug("body", zap.String("body", s))
	data := make([]NodeState, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		l.Error("Unable to unmarshal the checker response", zap.Error(err))
		return nil, err
	}

	check.successDetails = ""
	for _, item := range data {
		if item.Value != "UP" {
			err = fmt.Errorf("At least one node is DOWN")
			l.Error("Scylla checker error", zap.Error(err))
			return nil, err
		}
		check.successDetails += item.Key + "->UP..."
	}

	return check.successDetails, err
}

// NewScyllaChecker creates a health check for scylladb, as deployed in the youwol platform
// The simplest way to check sylla health is to provide the service URL
func NewScyllaChecker(cfg *checks.HTTPCheckConfig) *health.Config {

	sc, err := NewScyllaCheckerPrivate(cfg)
	if err != nil {
		return nil
	}
	return &health.Config{
		Check:            sc,
		ExecutionPeriod:  time.Duration(5) * time.Second,
		InitialDelay:     0,
		InitiallyPassing: false,
	}
}

// NewScyllaCheckerPrivate creates a new http check defined by the given config
func NewScyllaCheckerPrivate(config *checks.HTTPCheckConfig) (check checks.Check, err error) {
	l := utils.NewLogger()
	if config.URL == "" {
		err = fmt.Errorf("URL must not be empty")
		l.Error("Scylla checker initialization - Empty URL error", zap.Error(err))
		return nil, err
	}
	_, err = url.Parse(config.URL)
	if err != nil {
		err = errors.WithStack(err)
		l.Error("Scylla checker initialization - URL parsing", zap.Error(err))
		return nil, err
	}
	if config.CheckName == "" {
		err = fmt.Errorf("CheckName must not be empty")
		l.Error("Scylla checker initialization - Empty checker name error", zap.Error(err))
		return nil, err
	}

	if config.ExpectedStatus == 0 {
		config.ExpectedStatus = http.StatusOK
	}
	if config.Method == "" {
		config.Method = http.MethodGet
	}
	if config.Timeout == 0 {
		config.Timeout = time.Second
	}
	if config.Client == nil {
		config.Client = &http.Client{}
	}
	config.Client.Timeout = config.Timeout

	check = &scyllaChecker{
		config:         config,
		successDetails: fmt.Sprintf("Wait and see...%s", config.URL),
	}

	return check, nil
}
