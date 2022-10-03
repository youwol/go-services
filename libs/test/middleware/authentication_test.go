package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "gitlab.com/youwol/platform/libs/go-libs/middleware"
	"gitlab.com/youwol/platform/libs/go-libs/utils"
)

func TestAuthenticationMiddleware(t *testing.T) {

	h := middleware.NewAuthenticationMiddleware("test", utils.Local, nil)
	req := httptest.NewRequest("GET", "/toto", nil)
	recorder := httptest.NewRecorder() // Second server to provide the /ready endpoint
	h.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("No authorization should be granted without token (%d)", recorder.Code)
	}
}

// // Requires admin auth to work
// func TestGetGroups(t *testing.T) {
// 	middleware.NewAuthenticationMiddleware(nil) // Required to initialize the cache

// 	grp := middleware.GetGroupByPath(context.Background(), "/youwol-users/petronas")
// 	if grp == nil {
// 		t.Errorf("Group was not found")
// 	}
// }
