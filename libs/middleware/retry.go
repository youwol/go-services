package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"gitlab.com/youwol/platform/libs/go-libs/utils"
	"go.uber.org/zap"
)

// NewRetryMiddleware retries every request for which status is >= 400
func NewRetryMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		operation := func() error {
			next.ServeHTTP(rec, r)

			// Retry if error
			if rec.Code >= http.StatusBadRequest {
				return fmt.Errorf("%s", rec.Body.Bytes())
			}

			for k, v := range rec.HeaderMap {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
			return nil
		}

		b := &backoff.ExponentialBackOff{
			InitialInterval:     500 * time.Millisecond, // backoff.DefaultInitialInterval,
			RandomizationFactor: backoff.DefaultRandomizationFactor,
			Multiplier:          backoff.DefaultMultiplier,
			MaxInterval:         2 * time.Second,  // backoff.DefaultMaxInterval,
			MaxElapsedTime:      30 * time.Second, // backoff.DefaultMaxElapsedTime,
			Stop:                backoff.Stop,
			Clock:               backoff.SystemClock,
		}
		b.Reset()

		b1 := backoff.WithMaxRetries(b, 9)
		bo := backoff.WithContext(b1, r.Context())
		err := backoff.Retry(operation, bo)
		if err != nil {
			logger := utils.ContextLogger(r.Context())
			logger.Error("Retry backoff failed after several retries", zap.Error(err))
			next.ServeHTTP(w, r) // Provide the real response
		}
	})
}
