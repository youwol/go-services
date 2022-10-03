package handlers

import (
	"fmt"
	"platform/services/docdb/models"
	"platform/services/docdb/restapi/operations/index"

	"github.com/go-openapi/runtime/middleware"
	zap "go.uber.org/zap"
)

// AddIndexHandler handles the POST table request
func AddIndexHandler(params index.AddIndexParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)
	defer InvalidateSession(ctx, ksName)

	if err != nil {
		return index.NewAddIndexBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Let's build the query
	q := "CREATE INDEX "
	if params.Index.Name != nil {
		q += *params.Index.Name + " "
	}
	q += "ON " + params.TableName
	if len(params.Index.Identifier.PartitionKey) > 0 {
		q += "((" + params.Index.Identifier.PartitionKey
		q += "," + OwnerIDColumn + "," + OwnerNameColumn + "," + OwnerKindColumn + "),"
	} else {
		q += "("
	}
	if len(params.Index.Identifier.Option) > 0 {
		q += params.Index.Identifier.Option + "(" + *params.Index.Identifier.ColumnName + ")"
	} else {
		q += *params.Index.Identifier.ColumnName
	}
	q += ")"
	logger.Debug("built index query", zap.String("query", q))
	err = session.Query(q).Exec()
	if err != nil {
		logger.Error("Could not create new index", zap.Error(err), zap.Any("params", params), zap.String("query", q))
		return index.NewAddIndexBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return index.NewAddIndexCreated().WithPayload(&models.APIResponse{Message: "Index created"})
}

// DeleteIndexHandler handles the DELETE table request
func DeleteIndexHandler(params index.DeleteIndexParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)
	if err != nil {
		return index.NewDeleteIndexBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	defer InvalidateSession(ctx, ksName)

	// first check if the Keyspace exists
	q := fmt.Sprintf(`DROP INDEX IF EXISTS %s`, params.IndexName)
	err = session.Query(q).Exec()

	if err != nil {
		logger.Error("Could not delete existing table", zap.Error(err), zap.Any("params", params))
		return index.NewDeleteIndexBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	return index.NewDeleteIndexOK().WithPayload(&models.APIResponse{Message: "Index deleted"})
}
