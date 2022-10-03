package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"platform/services/storage/models"
	"strings"
	"time"

	minio "github.com/minio/minio-go"
	"github.com/patrickmn/go-cache"
	authz "gitlab.com/youwol/platform/libs/go-libs/middleware"
	utils "gitlab.com/youwol/platform/libs/go-libs/utils"
	zap "go.uber.org/zap"
)

var endpoint = os.Getenv("MINIO_HOST_PORT")
var accessKeyID = os.Getenv("MINIO_ACCESS_KEY")
var secretAccessKey = os.Getenv("MINIO_ACCESS_SECRET")
var c = cache.New(10*time.Minute, 10*time.Minute+1*time.Second)
var minioClientKey = "minio_client_" + os.Getenv("ENVIRONMENT")

// InitializeContext is used to set static data at server startup
// Those settings can be retrieved from command line or environment
// func InitializeContext(minioHost string) {
// 	endpoint = minioHost
// }

// GetHandlerContext retrieves the main objects that are necessary for handlers implementation
func GetHandlerContext(r *http.Request) (context.Context, *minio.Client, zap.Logger, error) {
	// Common initialization of request context
	ctx := r.Context()
	logger := utils.ContextLogger(ctx)
	client, err := getMinioClient(ctx)
	if err != nil {
		logger.Error("Could not create handler context", zap.Error(err))
	}
	return ctx, client, logger, err
}

func getMinioClient(ctx context.Context) (*minio.Client, error) {
	var err error
	client, found := c.Get(minioClientKey)
	if !found {
		client, err = minio.New(endpoint, accessKeyID, secretAccessKey, false)
		logger := utils.ContextLogger(ctx)
		if err == nil {
			logger.Info("Minio client initialized", zap.String("URL", endpoint))
			c.Set(minioClientKey, client, cache.DefaultExpiration)
		} else {
			logger.Error("Could not initialize minio client", zap.Error(err), zap.String("URL", endpoint))
		}
	}

	return client.(*minio.Client), err
}

// InvalidateClient empties the client cache and frees the current cached client
func InvalidateClient() {
	c.Delete(minioClientKey)
}

// IsAuthorized returns a boolean indicating if the current user is authorized to access the current object
func IsAuthorized(ctx context.Context, stats *minio.ObjectInfo, auth *models.Principal) bool {
	logger := utils.ContextLogger(ctx)
	ownerName := stats.Metadata.Get("X-Amz-Meta-Owner_name")
	ownerKind := stats.Metadata.Get("X-Amz-Meta-Owner_kind")
	owner := models.Owner{
		ID:   stats.Metadata.Get("X-Amz-Meta-Owner_id"),
		Name: &ownerName,
		Kind: &ownerKind,
	}
	logger.Debug("Owner retrieved", zap.Any("owner", owner))
	if ownerKind == "user" {
		return owner.ID == *auth.Sub
	} else if ownerKind == "group" {
		for _, group := range auth.MemberOf {
			if group.Path == ownerName || strings.HasPrefix(group.Path, ownerName) {
				return true
			}
		}
	}

	return false
}

// SetCredentials sets ownership metadata when uploading objects
func SetCredentials(ctx context.Context, owner *string, options *minio.PutObjectOptions, auth *models.Principal) error {
	if owner != nil && len(*owner) > 0 {
		// Check that we are member of this group
		bMatched := false
		for _, authGroup := range auth.MemberOf {
			if authGroup.Path == *owner { // Group name matches exactly
				options.UserMetadata["owner_id"] = authGroup.ID
				options.UserMetadata["owner_name"] = authGroup.Path
				options.UserMetadata["owner_kind"] = "group"
				bMatched = true
			} else if strings.HasPrefix(authGroup.Path, *owner) { // We are using a parent group, then we need to retrieve its information
				grp := authz.GetGroupByPath(ctx, "storage", *owner) // Get the parent group data
				if grp != nil {
					options.UserMetadata["owner_id"] = *grp.ID
					options.UserMetadata["owner_name"] = *grp.Path
					options.UserMetadata["owner_kind"] = "group"
					bMatched = true
				} else {
					return fmt.Errorf("Error retrieving parent group")
				}
			}
		}
		if !bMatched {
			return fmt.Errorf("You are not member of this group")
		}
	} else {
		// Set default ownership to the current user
		options.UserMetadata["owner_id"] = *auth.Sub
		options.UserMetadata["owner_name"] = auth.PreferredUsername
		options.UserMetadata["owner_kind"] = "user"
	}

	return nil
}

// GetIsolatedObjName return the target object name when isolated mode is ON
func GetIsolatedObjName(objName string, owner *string, auth *models.Principal, isolated *bool) string {
	var iso bool = true
	if isolated != nil {
		iso = *isolated
	}
	if !iso {
		return objName
	}

	prefix := auth.PreferredUsername
	if owner != nil && len(*owner) > 0 {
		prefix = *owner
	}
	prefix = strings.TrimPrefix(prefix, "/")
	prefix = strings.TrimSuffix(prefix, "/")
	return string(prefix + "/" + objName)
}
