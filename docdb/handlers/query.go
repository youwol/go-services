package handlers

import (
	"platform/services/docdb/models"
	"platform/services/docdb/restapi/operations/query"

	"github.com/go-openapi/runtime/middleware"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	zap "go.uber.org/zap"
)

// SelectQueryHandler handles the POST SELECT request
func SelectQueryHandler(params query.SelectQueryParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return query.NewSelectQueryBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Costly to do it every time, should be cached
	_, cols, partKey, _, err := getMeta(ctx, session, ksName, params.TableName)
	if err != nil {
		logger.Error("Error retrieving table metadata", zap.Error(err))
		return query.NewSelectQueryBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ownerID, ownerName, ownerKind, err := GetCredentials(ctx, params.Owner, auth)
	if err != nil {
		return query.NewSelectQueryUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if params.Select.Query == nil {
		params.Select.Query = &models.Query{}
	}
	if params.Select.Query.WhereClause == nil {
		params.Select.Query.WhereClause = make([]*models.QueryRelation, 0)
	}

	// append authz condiditons only if the query is made with partition key,
	// else it will not work on indexes
	partKeyQuery := false
	for _, wcp := range params.Select.Query.WhereClause {
		for _, partKeyMeta := range partKey {
			if partKeyMeta.Name == *wcp.Column {
				partKeyQuery = true
				break
			}
		}
	}
	authzChecked := partKeyQuery || (params.Select.AllowFiltering != nil && *params.Select.AllowFiltering)
	if authzChecked {
		params.Select.Query.WhereClause = append(params.Select.Query.WhereClause,
			&models.QueryRelation{Column: &OwnerIDColumn, Relation: models.RelationOperatorEq, Term: ownerID},
			&models.QueryRelation{Column: &OwnerKindColumn, Relation: models.RelationOperatorEq, Term: ownerKind},
			&models.QueryRelation{Column: &OwnerNameColumn, Relation: models.RelationOperatorEq, Term: ownerName},
		)
	}

	// params.Select.SelectClause
	selectBuilder := qb.Select(params.TableName).Limit(uint(*params.Select.MaxResults))
	buildSelectClause(selectBuilder, &params.Select.SelectClause)
	whereValues := make(qb.M)
	if params.Select.Query != nil {
		buildWhereClause(selectBuilder, nil, &whereValues, &params.Select.Query.WhereClause)
		buildOrderByClause(selectBuilder, &params.Select.Query.OrderingClause)
		if params.Select.AllowFiltering != nil && *params.Select.AllowFiltering == true {
			selectBuilder.AllowFiltering() // TODO: check when to use
		}
	}
	stmt, names := selectBuilder.ToCql()

	// allocating the result structure, unfortunately we cannot know how many records in advance
	rowResult := make(map[string]interface{})
	ret := &models.SelectResponse{}
	if params.Select.Mode == nil || *params.Select.Mode == models.SelectStatementModeDocuments {
		ret.Documents = make([]map[string]interface{}, 0)
	} else {
		ret.Columns = make(map[string][]interface{})
	}

	q := gocqlx.Query(session.Query(stmt), names).BindMap(whereValues)
	iter := q.Iter()
	rowIndex := 0
	for iter.MapScan(rowResult) {
		// Append results
		cqlToJSON(&rowResult, &cols)
		if !authzChecked { // check auhtz on the record
			if rowResult[OwnerIDColumn].(string) != ownerID {
				continue
			}
		}
		if params.Select.Mode == nil || *params.Select.Mode == models.SelectStatementModeDocuments {
			ret.Documents = append(ret.Documents.([]map[string]interface{}), rowResult)
		} else {
			// Transform each row map and append it to a map of arrays
			for k, v := range rowResult {
				cols := ret.Columns.(map[string][]interface{})
				_, found := cols[k]
				if !found {
					cols[k] = make([]interface{}, 0)
				}
				cols[k] = append(cols[k], v)
			}
		}
		rowResult = make(map[string]interface{})
		rowIndex++ // TODO: check if max results have been reached and return an iterator
	}
	defer q.Release()

	return query.NewSelectQueryOK().WithPayload(ret)
}

// DeleteQueryHandler handles the DELETE query request
func DeleteQueryHandler(params query.DeleteQueryParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	_, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return query.NewDeleteQueryBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ownerID, ownerName, ownerKind, err := GetCredentials(ctx, params.Owner, auth)
	if err != nil {
		return query.NewSelectQueryUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	params.Delete.WhereClause = append(params.Delete.WhereClause,
		&models.QueryRelation{Column: &OwnerIDColumn, Relation: models.RelationOperatorEq, Term: ownerID},
		&models.QueryRelation{Column: &OwnerKindColumn, Relation: models.RelationOperatorEq, Term: ownerKind},
		&models.QueryRelation{Column: &OwnerNameColumn, Relation: models.RelationOperatorEq, Term: ownerName},
	)

	// Build the query
	deleteBuilder := qb.Delete(params.TableName)
	deleteBuilder.Columns(params.Delete.SimpleSelection...)
	// deleteBuilder.AllowFiltering() // TODO: check when to use
	whereValues := make(qb.M)
	buildWhereClause(nil, deleteBuilder, &whereValues, &params.Delete.WhereClause)
	stmt, names := deleteBuilder.ToCql()

	// Run the query
	q := gocqlx.Query(session.Query(stmt), names).BindMap(whereValues)
	if err := q.ExecRelease(); err != nil {
		defer InvalidateSession(ctx, params.KeyspaceName)
		logger.Error("Issue while trying to DELETE entities", zap.Error(err), zap.Any("params", params))
		return query.NewDeleteQueryBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return query.NewDeleteQueryOK().WithPayload(&models.APIResponse{Message: "Documents deleted"})
}
