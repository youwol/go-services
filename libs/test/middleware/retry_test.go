package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	// "time"

	middleware "gitlab.com/youwol/platform/libs/go-libs/middleware"
)

func TestRetryMiddlewareOk(t *testing.T) {
	// Create middleware over a Ok function
	h := middleware.NewRetryMiddleware(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Send response to be tested
			rw.WriteHeader(http.StatusOK)
			io.WriteString(rw, "OK")
		}),
	)

	// check the response and response time
	rq := httptest.NewRequest("GET", "http://toto.com", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, rq)

	if w.Code != http.StatusOK {
		t.Errorf("Retry middleware failed for OK response")
	}

	if w.Body.String() != "OK" {
		t.Errorf("Retry middleware should not modify the response (ret=%s)", w.Body.String())
	}

}

func TestRetryMiddlewareNOTOk(t *testing.T) {
	// Create middleware over a Ok function
	h := middleware.NewRetryMiddleware(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Send response to be tested
			rw.WriteHeader(http.StatusNotFound)
			io.WriteString(rw, "FAILED")
		}),
	)

	// check the response and response time
	rq := httptest.NewRequest("GET", "http://toto.com", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, rq)

	if w.Code != http.StatusNotFound {
		t.Errorf("Retry middleware should have failed for a NOT OK response")
	}
	if w.Body.String() != "FAILED" {
		t.Errorf("Retry middleware should have returned the handler body (ret=%s)", w.Body.String())
	}
}

func TestRetryMiddlewareOneFailedOneOk(t *testing.T) {
	bFail := true
	// Create middleware over a Ok function
	h := middleware.NewRetryMiddleware(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Send response to be tested
			if bFail {
				rw.WriteHeader(http.StatusNotFound)
				io.WriteString(rw, "FAILED")
				bFail = false
			} else {
				rw.WriteHeader(http.StatusOK)
				io.WriteString(rw, "OK")
			}
		}),
	)

	// check the response and response time
	rq := httptest.NewRequest("GET", "http://toto.com", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, rq)

	if w.Code != http.StatusOK {
		t.Errorf("Retry middleware should not have failed after retry")
	}
	if w.Body.String() != "OK" {
		t.Errorf("Retry middleware should have returned the correct body (ret=%s)", w.Body.String())
	}
}
