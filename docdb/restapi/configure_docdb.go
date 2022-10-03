// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strings"
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.com/patrickmn/go-cache"
	healthchecks "gitlab.com/youwol/platform/libs/go-libs/health"
	healthcheckers "gitlab.com/youwol/platform/libs/go-libs/health/checkers"
	middleware "gitlab.com/youwol/platform/libs/go-libs/middleware"
	utils "gitlab.com/youwol/platform/libs/go-libs/utils"

	"platform/services/docdb/handlers"
	"platform/services/docdb/models"
	"platform/services/docdb/restapi/operations"
	"platform/services/docdb/restapi/operations/document"
	"platform/services/docdb/restapi/operations/index"
	"platform/services/docdb/restapi/operations/keyspace"
	"platform/services/docdb/restapi/operations/query"
	"platform/services/docdb/restapi/operations/table"
)

var tokenCache *cache.Cache

//go:generate swagger generate server --target ..\..\docdb --name Docdb --spec ..\api\docdb-api.yaml
var customFlags = struct {
	NoCheckers bool `long:"nocheckers" description:"Disable liveness and readiness probes"`
}{}

func configureFlags(api *operations.DocdbAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Additional flags",
			LongDescription:  "",
			Options:          &customFlags,
		},
	}
}

func checkEnv(logger *zap.Logger) {
	envs := []string{"SCYLLA_HOSTS", "KEYCLOAK_LOGIN", "KEYCLOAK_PASSWORD", "KEYCLOAK_HOST",
		"REDIS_HOST_PORT", "ENVIRONMENT"}
	for _, e := range envs {
		if os.Getenv(e) == "" {
			logger.Fatal("Environment variable not set", zap.String("Variable", e))
		}
	}
	logger.Info("Authorization", zap.String("KEYCLOAK_HOST", os.Getenv("KEYCLOAK_HOST")))
}

func configureAPI(api *operations.DocdbAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf
	api.Logger = utils.APILogger
	checkEnv(utils.NewLogger())
	tokenCache = cache.New(time.Hour, time.Hour+time.Second)

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Static initialization
	handlers.InitConverters()

	// Applies when the "Authorization" header is set
	api.BearerAuth = func(token string) (*models.Principal, error) {
		var err error
		token = strings.TrimPrefix(token, "Bearer ")
		pcp, bFound := tokenCache.Get(token)
		principal := models.Principal{}
		if !bFound {
			// RE-FACTOR ME: there are two calls to GetUserInfo for each request!!
			mapUser := middleware.GetUserInfo(context.Background(), "docdb", token)
			err = mapstructure.Decode(mapUser, &principal)
			tokenCache.Set(token, principal, time.Hour)
		} else {
			principal = pcp.(models.Principal)
		}
		return &principal, err
	}

	// plug our handlers onto the API...

	// Keyspace endpoints
	api.KeyspaceAddKeyspaceHandler = keyspace.AddKeyspaceHandlerFunc(handlers.AddKeyspaceHandler)
	api.KeyspaceDeleteKeyspaceHandler = keyspace.DeleteKeyspaceHandlerFunc(handlers.DeleteKeyspaceHandler)
	api.KeyspaceGetKeyspaceHandler = keyspace.GetKeyspaceHandlerFunc(handlers.GetKeyspaceHandler)
	api.KeyspaceUpdateKeyspaceHandler = keyspace.UpdateKeyspaceHandlerFunc(handlers.UpdateKeyspaceHandler)
	api.KeyspaceGetKeyspacesHandler = keyspace.GetKeyspacesHandlerFunc(handlers.GetKeyspacesHandler)

	// Table endpoints
	api.TableAddTableHandler = table.AddTableHandlerFunc(handlers.AddTableHandler)
	api.TableDeleteTableHandler = table.DeleteTableHandlerFunc(handlers.DeleteTableHandler)
	api.TableUpdateTableHandler = table.UpdateTableHandlerFunc(handlers.UpdateTableHandler)
	api.TableGetTableHandler = table.GetTableHandlerFunc(handlers.GetTableHandler)
	api.TableGetTablesHandler = table.GetTablesHandlerFunc(handlers.GetTablesHandler)

	// Index endpoints
	api.IndexAddIndexHandler = index.AddIndexHandlerFunc(handlers.AddIndexHandler)
	api.IndexDeleteIndexHandler = index.DeleteIndexHandlerFunc(handlers.DeleteIndexHandler)

	// Document endpoints
	api.DocumentAddDocumentHandler = document.AddDocumentHandlerFunc(handlers.AddDocumentHandler)
	api.DocumentGetDocumentHandler = document.GetDocumentHandlerFunc(handlers.GetDocumentHandler)
	api.DocumentDeleteDocumentHandler = document.DeleteDocumentHandlerFunc(handlers.DeleteDocumentHandler)
	api.DocumentUpdateDocumentHandler = document.UpdateDocumentHandlerFunc(handlers.UpdateDocumentHandler)

	// Query endpoints
	api.QuerySelectQueryHandler = query.SelectQueryHandlerFunc(handlers.SelectQueryHandler)
	api.QueryDeleteQueryHandler = query.DeleteQueryHandlerFunc(handlers.DeleteQueryHandler)

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
	// Middlewares are run in the reverse order of their installation
	// handler = gziphandler.GzipHandler(handler) looks slower than without gzipping
	handler = middleware.NewAuthenticationMiddleware("docdb", utils.Redis, handler)
	// handler = middleware.NewRetryMiddleware(handler)
	handler = middleware.NewContextLoggerMiddleware(handler)

	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	// Create readiness and liveness checkers
	if !customFlags.NoCheckers {
		// Configure a checker for each node
		hosts := strings.Split(os.Getenv("SCYLLA_HOSTS"), ",")
		h := hosts[0] // No need to check all the nodes only the first seed is enough
		livenessCheckers := []*health.Config{
			healthcheckers.NewOkChecker("I am alive"),
			healthcheckers.NewScyllaChecker(
				&checks.HTTPCheckConfig{
					CheckName: "scylla-checker",
					URL:       "http://" + h + ":10000/failure_detector/simple_states",
				},
			),
		}
		handler = middleware.NewHealthMiddleware(handler, livenessCheckers, healthchecks.ProbeTypeLiveness)

		readinessCheckers := []*health.Config{
			healthcheckers.NewOkChecker("I am ready"),
			healthcheckers.NewScyllaChecker(
				&checks.HTTPCheckConfig{
					CheckName: "scylla-checker",
					URL:       "http://" + h + ":10000/failure_detector/simple_states",
				}),
		}
		handler = middleware.NewHealthMiddleware(handler, readinessCheckers, healthchecks.ProbeTypeReadiness)
	}

	return handler
}
