package middleware_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "gitlab.com/youwol/platform/libs/go-libs/middleware"
)

func BenchmarkWithoutLogger(t *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	for i := 0; i < t.N; i++ {
		_, err := http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkWithLogger(t *testing.B) {
	loggerFunc := middleware.NewContextLoggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	ts := httptest.NewServer(loggerFunc)
	defer ts.Close()
	for i := 0; i < t.N; i++ {
		_, err := http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}
	}
}
