package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/youwol/platform/libs/go-libs/utils"

	"github.com/Nerzal/gocloak/v5"
	"github.com/fatih/structs"
	"go.uber.org/zap"
)

var keycloakHost = os.Getenv("KEYCLOAK_HOST")
var keycloakClientId = os.Getenv("OPENID_CLIENT_ID")
var keycloakClientSecret = os.Getenv("OPENID_CLIENT_SECRET")

// NewAuthenticationMiddleware fetches user info from the bearer token
// It has to be installed after the ContextLoggerMiddleware
func NewAuthenticationMiddleware(cacheName string, ct utils.CacheType, next http.Handler) http.Handler {
	utils.InitCache(cacheName, ct, 1*time.Hour)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ready" || r.URL.Path == "/alive" {
			next.ServeHTTP(w, r)
		} else {
			ctx := r.Context()
			logger := utils.ContextLogger(ctx)
			token := r.Header.Get("Authorization")
			// Remove the bearer part to leave only the token
			token = strings.TrimPrefix(token, "Bearer ")
			if len(token) == 0 {
				logger.Error("No authorization header was found", zap.Any("Headers", r.Header))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ret := GetUserInfo(r.Context(), cacheName, token)

			// serve next middleware if auth is ok
			if ret != nil {
				next.ServeHTTP(w, r)
			} else {
				logger.Error("Authorization failure")
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	})
}

// GetUserInfo fetches user information from the access token (either from cache or keycloak)
func GetUserInfo(ctx context.Context, cacheName string, token string) map[string]interface{} {
	userMap := make(map[string]interface{})
	bFound := utils.GetCache(ctx, cacheName, token, &userMap)
	if !bFound { // Fetch user info from keycloak and store it in cache
		logger := utils.ContextLogger(ctx)

		client := gocloak.NewClient(keycloakHost)
		//JWT, err := client.LoginAdmin(keycloakLogin, keycloakPassword, "master")
		_, err := client.LoginClient(keycloakClientId, keycloakClientSecret, "youwol")
		if err != nil {
			logger.Error("Could not login into keycloak", zap.Error(err))
			return nil
		}
		userinfo, err := client.GetUserInfo(token, "youwol")
		if err != nil {
			logger.Error("Could not retrieve user info", zap.Error(err))
			return nil
		}
		logger.Info("User info retrieved", zap.Any("UserInfo", *userinfo))
		//userGroups, err := client.GetUserGroups(JWT.AccessToken, "youwol", *userinfo.Sub)
		//if err != nil {
		//	logger.Error("Could not retrieve current users groups", zap.Error(err))
		//	return nil
		//}

		groupID := "b289d90d-786e-49b9-a2d2-6ad76ab73ba5"
		groupPath := "/youwol-users"

		youwolUsersGroup := [1]*gocloak.Group{{
			ID:   &groupID,
			Path: &groupPath,
		}} // We use the access token as a key to store user info
		userMap = structs.Map(*userinfo)
		userMap["MemberOf"] = youwolUsersGroup

		err = utils.SetCache(ctx, cacheName, token, userMap)
		if err != nil {
			logger.Error("Could not write to cache", zap.Error(err))
			return nil
		}
	}

	return userMap
}

// GetGroupByPath finds group information from its path and from the access token
func GetGroupByPath(ctx context.Context, cacheName string, groupPath string) *gocloak.Group {
	group := &gocloak.Group{}
	bFound := utils.GetCache(ctx, cacheName, groupPath, &group)
	if !bFound { // Fetch user info from keycloak and store it in cache
		logger := utils.ContextLogger(ctx)
		groupHierarchy := strings.Split(groupPath, "/")
		if len(groupHierarchy) < 1 {
			logger.Error("Could not parse group path")
			return nil
		}

		client := gocloak.NewClient(keycloakHost)
		//JWT, err := client.LoginAdmin(keycloakLogin, keycloakPassword, "master")
		JWT, err := client.LoginClient(keycloakClientId, keycloakClientSecret, "master")
		if err != nil {
			logger.Error("Could not login into keycloak", zap.Error(err))
			return nil
		}
		groups, err := client.GetGroups(JWT.AccessToken, "youwol", gocloak.GetGroupsParams{Search: &groupHierarchy[1]})
		if err != nil {
			logger.Error("Could not retrieve group hierarchy", zap.Error(err))
			return nil
		}
		logger.Info("Root group retrieved", zap.Any("RootInfo", groups))

		// We use the access token as a key to store user info
		if len(groups) == 1 {
			grp, err := findGroupPathInHierarchy(groups[0], groupPath)
			if err != nil {
				logger.Error("Error searching for group", zap.Error(err))
				return nil
			}
			// group = structs.Map(groups[0])
			group = grp
			utils.SetCache(ctx, cacheName, groupPath, grp)
		} else {
			logger.Error("There were multiple group matches")
			return nil
		}
	}

	return group
}

// Recursive function to find the target group in a group hierarchy
func findGroupPathInHierarchy(hierarchy *gocloak.Group, groupPath string) (*gocloak.Group, error) {
	if hierarchy == nil {
		return nil, nil
	}

	if *hierarchy.Path == groupPath {
		return hierarchy, nil
	}

	for _, child := range hierarchy.SubGroups {
		res, err := findGroupPathInHierarchy(child, groupPath)
		if res != nil && err == nil {
			return res, err
		}
	}

	return nil, fmt.Errorf("Could not find target group path")
}
