package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"gitlab.com/youwol/platform/libs/go-libs/utils"

	"go.uber.org/zap"
)

// NewContextLoggerMiddleware creates a request context
// It also logs requests and responses
func NewContextLoggerMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture the request body
		reqBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

		// build request context
		r = r.WithContext(utils.CreateRequestContext(*r))
		logger := utils.ContextLogger(r.Context())

		// record response writing
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		// copy everything from response recorder
		// to actual response writer
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		// body := rec.Body.Bytes()
		w.Write(rec.Body.Bytes())

		end := time.Now()

		// log response : do not log encoded responses bodies
		defer func() {
			ct := rec.HeaderMap.Get("Content-Type")
			ce := rec.HeaderMap.Get("Content-Encoding")
			data := make(map[string]interface{})
			err := json.Unmarshal(rec.Body.Bytes(), &data)
			if ct == "application/json" && ce != "gzip" && len(rec.Body.Bytes()) > 0 && len(rec.Body.Bytes()) < 1024 && err == nil {
				logger.Info("Response sent", zap.Int("status", rec.Code), zap.Any("body", data), zap.String("duration", fmt.Sprintf("%vms", end.Sub(start).Nanoseconds()/1e6))) // pre-log
			} else {
				logger.Info("Response sent", zap.Int("status", rec.Code), zap.String("duration", fmt.Sprintf("%vms", end.Sub(start).Nanoseconds()/1e6))) // pre-log
			}
		}()

		// log request
		defer func() {
			ct := r.Header.Get("Content-Type")
			ce := r.Header.Get("Content-Encoding")
			data := make(map[string]interface{})
			err := json.Unmarshal(reqBody, &data)
			if ct == "application/json" && ce != "gzip" && len(reqBody) > 0 && err == nil {
				logger.Info("Request received", zap.Any("body", data), zap.String("path", r.RequestURI), zap.String("method", r.Method))
			} else {
				logger.Info("Request received", zap.String("path", r.RequestURI), zap.String("method", r.Method))
			}
		}()
	})
}
