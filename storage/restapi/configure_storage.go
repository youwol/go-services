// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strings"
	"time"

	"platform/services/storage/handlers"
	"platform/services/storage/models"
	"platform/services/storage/restapi/operations"
	"platform/services/storage/restapi/operations/bucket"
	"platform/services/storage/restapi/operations/file"
	"platform/services/storage/restapi/operations/object"
	"platform/services/storage/restapi/operations/objectinfo"
	"platform/services/storage/restapi/operations/objects"

	health "github.com/AppsFlyer/go-sundheit"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/mitchellh/mapstructure"
	"github.com/patrickmn/go-cache"
	healthchecks "gitlab.com/youwol/platform/libs/go-libs/health"
	healthcheckers "gitlab.com/youwol/platform/libs/go-libs/health/checkers"
	"gitlab.com/youwol/platform/libs/go-libs/middleware"
	utils "gitlab.com/youwol/platform/libs/go-libs/utils"
	"go.uber.org/zap"
)

var tokenCache *cache.Cache

//go:generate swagger generate server --target ..\..\storage --name Storage --spec ..\api\storage-api.yaml

var customFlags = struct {
	NoCheckers bool `long:"nocheckers" description:"Disable liveness and readiness probes"`
	NoAuth     bool `long:"noauth" description:"Disable authentication middleware"`
}{}

func configureFlags(api *operations.StorageAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "Additional flags",
			LongDescription:  "",
			Options:          &customFlags,
		},
	}
}

func checkEnv(logger *zap.Logger) {
	envs := []string{"MINIO_HOST_PORT", "MINIO_ACCESS_KEY", "MINIO_ACCESS_SECRET",
		"REDIS_HOST_PORT", "KEYCLOAK_HOST", "OPENID_CLIENT_ID", "OPENID_CLIENT_SECRET"}
	for _, e := range envs {
		if os.Getenv(e) == "" {
			logger.Fatal("Environment variable not set", zap.String("Variable", e))
		}
	}
	logger.Info("Authorization", zap.String("KEYCLOAK_HOST", os.Getenv("KEYCLOAK_HOST")))
}

func configureAPI(api *operations.StorageAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = utils.APILogger
	tokenCache = cache.New(time.Hour, time.Hour+time.Second)
	checkEnv(utils.NewLogger())

	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer
	api.JSONProducer = runtime.JSONProducer()
	api.BinProducer = runtime.ByteStreamProducer()

	// Applies when the "Authorization" header is set
	api.BearerAuth = func(token string) (*models.Principal, error) {
		var err error
		principal := models.Principal{}
		if !customFlags.NoAuth {
			token = strings.TrimPrefix(token, "Bearer ")
			pcp, bFound := tokenCache.Get(token)
			if !bFound {
				// RE-FACTOR ME: there are two calls to GetUserInfo for each request!!
				mapUser := middleware.GetUserInfo(context.Background(), "storage", token)
				err = mapstructure.Decode(mapUser, &principal)
				tokenCache.Set(token, principal, time.Hour)
			} else {
				principal = pcp.(models.Principal)
			}
		}
		return &principal, err
	}

	api.BucketAddBucketHandler = bucket.AddBucketHandlerFunc(handlers.AddBucketHandler)
	api.BucketDeleteBucketHandler = bucket.DeleteBucketHandlerFunc(handlers.DeleteBucketHandler)
	api.BucketGetBucketsHandler = bucket.GetBucketsHandlerFunc(handlers.GetBucketsHandler)

	api.ObjectAddObjectHandler = object.AddObjectHandlerFunc(handlers.AddObjectHandler)
	api.ObjectGetObjectHandler = object.GetObjectHandlerFunc(handlers.GetObjectHandler)
	api.ObjectDeleteObjectHandler = object.DeleteObjectHandlerFunc(handlers.DeleteObjectHandler)

	api.ObjectinfoGetObjectInfoHandler = objectinfo.GetObjectInfoHandlerFunc(handlers.GetObjectInfoHandler)

	api.ObjectsGetObjectsHandler = objects.GetObjectsHandlerFunc(handlers.GetObjectsHandler)
	api.ObjectsDeleteObjectsHandler = objects.DeleteObjectsHandlerFunc(handlers.DeleteObjectsHandler)

	api.FileAddFileHandler = file.AddFileHandlerFunc(handlers.AddFileHandler)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	if !customFlags.NoAuth {
		handler = middleware.NewAuthenticationMiddleware("storage", utils.Redis, handler)
	}
	handler = middleware.NewContextLoggerMiddleware(handler)

	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	if !customFlags.NoCheckers {
		minioConfig := healthcheckers.NewMinioChecker(os.Getenv("MINIO_HOST_PORT"))

		livenessCheckers := []*health.Config{minioConfig}
		handler = middleware.NewHealthMiddleware(handler, livenessCheckers, healthchecks.ProbeTypeLiveness)

		readinessCheckers := []*health.Config{minioConfig}
		handler = middleware.NewHealthMiddleware(handler, readinessCheckers, healthchecks.ProbeTypeReadiness)
	}
	return handler
}
