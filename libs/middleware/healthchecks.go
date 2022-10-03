package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	health "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	healthchecks "gitlab.com/youwol/platform/libs/go-libs/health"
	"gitlab.com/youwol/platform/libs/go-libs/utils"
	"go.uber.org/zap"
)

// NewHealthMiddleware create a new handler for readiness probe
func NewHealthMiddleware(next http.Handler, checkers []*health.Config, pt healthchecks.ProbeType) http.Handler {

	h, err := healthchecks.NewHealthManager(checkers, pt)
	if err != nil {
		l := utils.NewLogger()
		l.Error("Could not create health check manager", zap.Error(err))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ready" && pt == healthchecks.ProbeTypeReadiness ||
			r.URL.Path == "/alive" && pt == healthchecks.ProbeTypeLiveness {
			results, healthy := (*h).Results()
			w.Header().Set("Content-Type", "application/json")
			if healthy {
				w.WriteHeader(200)
			} else {
				w.WriteHeader(503)
			}

			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "\t")
			var err error
			if r.URL.Query().Get("type") == healthhttp.ReportTypeShort {
				shortResults := make(map[string]string)
				for k, v := range results {
					if v.IsHealthy() {
						shortResults[k] = "PASS"
					} else {
						shortResults[k] = "FAIL"
					}
				}

				err = encoder.Encode(shortResults)
			} else {
				err = encoder.Encode(results)
			}

			if err != nil {
				_, _ = fmt.Fprintf(w, "Failed to render results JSON: %s", err)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
