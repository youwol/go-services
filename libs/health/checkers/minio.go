package healthchecks

import (
	"fmt"
	"os"
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/minio/minio-go"
	"gitlab.com/youwol/platform/libs/go-libs/utils"
	"go.uber.org/zap"
)

var endpoint = os.Getenv("MINIO_HOST_PORT")
var accessKeyID = os.Getenv("MINIO_ACCESS_KEY")
var secretAccessKey = os.Getenv("MINIO_ACCESS_SECRET")

type minioChecker struct {
	url            string
	successDetails string
}

func (check *minioChecker) Name() string {
	return "minio-checker"
}

func (check *minioChecker) Execute() (details interface{}, err error) {
	details = check.url

	// prepare the request
	l := utils.NewLogger()

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		l.Error("Unable to create minio client", zap.Error(err))
		return nil, err
	}

	_, err = client.ListBuckets()
	if err != nil {
		err = fmt.Errorf("Minio checker returned error")
		l.Error("Minio checker error", zap.Error(err))
		return details, err
	}

	check.successDetails = check.url + "->UP"

	return check.successDetails, err
}

// NewMinioChecker creates a health check for scylladb, as deployed in the youwol platform
// The simplest way to check sylla health is to provide the service URL
func NewMinioChecker(url string) *health.Config {

	sc, err := NewMinioCheckerPrivate(url)
	if err != nil {
		return nil
	}
	return &health.Config{
		Check:            sc,
		ExecutionPeriod:  time.Duration(2) * time.Second,
		InitialDelay:     0,
		InitiallyPassing: false,
	}
}

// NewMinioCheckerPrivate creates a new http check defined by the given config
func NewMinioCheckerPrivate(sURL string) (check checks.Check, err error) {
	l := utils.NewLogger()
	if sURL == "" {
		err = fmt.Errorf("URL must not be empty")
		l.Error("Minio checker initialization - Empty URL error", zap.Error(err))
		return nil, err
	}

	check = &minioChecker{
		url:            sURL,
		successDetails: fmt.Sprintf("Wait and see...%s", sURL),
	}

	return check, nil
}
