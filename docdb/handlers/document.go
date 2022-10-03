package handlers

import (
	"fmt"
	"platform/services/docdb/models"
	"platform/services/docdb/restapi/operations/document"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
	zap "go.uber.org/zap"
)

func setKeys(data *map[string]interface{}, part []*gocql.ColumnMetadata, cluster []*gocql.ColumnMetadata, partitionKey []string, clusteringKey []string) error {
	if len(part) != len(partitionKey) {
		return fmt.Errorf("Inconsistent partition key")
	}

	for i, v := range part {
		(*data)[v.Name] = partitionKey[i]
	}

	if len(cluster) != len(clusteringKey) {
		return fmt.Errorf("Inconsistent clustering key")
	}

	for i, v := range cluster {
		(*data)[v.Name] = clusteringKey[i]
	}

	return nil
}

// AddDocumentHandler handles the POST document request
func AddDocumentHandler(params document.AddDocumentParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return document.NewAddDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Costly to do it every time, should be cached
	m, cols, _, _, err := getMeta(ctx, session, ksName, params.TableName)
	if err != nil {
		logger.Error("Error retrieving table metadata", zap.Error(err))
		return document.NewAddDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ownerID, ownerName, ownerKind, err := GetCredentials(ctx, params.Owner, auth)
	if err != nil {
		return document.NewAddDocumentUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	doc := params.Document.(map[string]interface{})
	doc[OwnerIDColumn] = ownerID
	doc[OwnerNameColumn] = ownerName
	doc[OwnerKindColumn] = ownerKind
	err = jsonToCql(&doc, &cols)
	if err != nil {
		logger.Error("Error converting the input data", zap.Error(err))
		return document.NewAddDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// This automatically builds the insert query
	var documentTable = table.New(*m)
	stmt, names := documentTable.Insert()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(doc)
	if err := q.ExecRelease(); err != nil {
		defer InvalidateSession(ctx, ksName)
		logger.Error("Issue while trying to insert the entity", zap.Error(err), zap.Any("params", params), zap.Any("schema", m))
		return document.NewAddDocumentBadRequest().WithPayload(
			&models.APIResponse{
				Message: err.Error(),
			})
	}

	// build the entity ID
	entityID := make(map[string]interface{})
	for _, pk := range m.PartKey {
		if pk != OwnerIDColumn && pk != OwnerNameColumn && pk != OwnerKindColumn {
			entityID[pk] = doc[pk]
		}
	}
	return document.NewAddDocumentCreated().WithPayload(entityID)
}

// GetDocumentHandler handles the GET entity request
func GetDocumentHandler(params document.GetDocumentParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return document.NewGetDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Costly to do it every time, should be cached
	metaTable, columns, part, cluster, err := getMeta(ctx, session, ksName, params.TableName)
	if err != nil {
		logger.Error("Error retrieving table metadata", zap.Error(err))
		return document.NewGetDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ownerID, ownerName, ownerKind, err := GetCredentials(ctx, params.Owner, auth)
	if err != nil {
		return document.NewGetDocumentUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	params.PartitionKey = append(params.PartitionKey, ownerID, ownerName, ownerKind)

	in := make(map[string]interface{})
	err = setKeys(&in, part, cluster, params.PartitionKey, params.ClusteringKey)
	if err != nil {
		logger.Error("Mismatch with the document keys", zap.Error(err))
		return document.NewGetDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	err = jsonToCql(&in, &columns)
	if err != nil {
		logger.Error("Error converting the input data", zap.Error(err))
		return document.NewGetDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	out := make(map[string]interface{})

	// This automatically builds the SELECT * query
	var entityTable = table.New(*metaTable)
	stmt, names := entityTable.Get()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(in)
	defer q.Release()
	if err := q.Query.MapScan(out); err != nil {
		InvalidateSession(ctx, ksName)
		logger.Error("Issue while trying to GET the entity", zap.Error(err), zap.Any("params", params), zap.Any("schema", metaTable))
		return document.NewGetDocumentNotFound().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	err = cqlToJSON(&out, &columns)
	if err != nil {
		logger.Error("Error converting the output data", zap.Error(err))
		return document.NewGetDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return document.NewGetDocumentOK().WithPayload(&out)
}

// DeleteDocumentHandler handles the DELETE entity request
func DeleteDocumentHandler(params document.DeleteDocumentParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return document.NewDeleteDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Costly to do it every time, should be cached
	m, _, part, cluster, err := getMeta(ctx, session, ksName, params.TableName)
	if err != nil {
		logger.Error("Error retrieving table metadata", zap.Error(err))
	}

	ownerID, ownerName, ownerKind, err := GetCredentials(ctx, params.Owner, auth)
	if err != nil {
		return document.NewDeleteDocumentUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	params.PartitionKey = append(params.PartitionKey, ownerID, ownerName, ownerKind)

	in := make(map[string]interface{})
	err = setKeys(&in, part, cluster, params.PartitionKey, params.ClusteringKey)
	if err != nil {
		logger.Error("Mismatch with the document keys", zap.Error(err))
		return document.NewGetDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// This automatically builds the DELETE * query
	var entityTable = table.New(*m)
	stmt, names := entityTable.Delete()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(in)
	if err := q.ExecRelease(); err != nil {
		defer InvalidateSession(ctx, ksName)
		logger.Error("Issue while trying to GET the entity", zap.Error(err), zap.Any("params", params), zap.Any("schema", m))
		return document.NewDeleteDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return document.NewDeleteDocumentOK().WithPayload(&models.APIResponse{Message: "Document deleted"})
}

// UpdateDocumentHandler handles the POST entity request
func UpdateDocumentHandler(params document.UpdateDocumentParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return document.NewUpdateDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Costly to do it every time, should be cached
	m, cols, partition, cluster, err := getMeta(ctx, session, ksName, params.TableName)
	if err != nil {
		logger.Error("Error retrieving table metadata", zap.Error(err))
	}

	ownerID, ownerName, ownerKind, err := GetCredentials(ctx, params.Owner, auth)
	if err != nil {
		return document.NewUpdateDocumentUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	doc := params.Document.(map[string]interface{})
	doc[OwnerIDColumn] = ownerID
	doc[OwnerNameColumn] = ownerName
	doc[OwnerKindColumn] = ownerKind

	err = jsonToCql(&doc, &cols)
	if err != nil {
		logger.Error("Error converting the input data", zap.Error(err))
		return document.NewUpdateDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	columnsToUpdate := make([]string, 0)
	for k := range doc {
		var foundInKeys = false
		for _, col := range partition {
			if col.Name == k {
				foundInKeys = true
				break
			}
		}
		for _, col := range cluster {
			if col.Name == k {
				foundInKeys = true
				break
			}
		}
		if !foundInKeys {
			columnsToUpdate = append(columnsToUpdate, k)
		}
	}

	// This automatically builds the insert query
	var entityTable = table.New(*m)
	stmt, names := entityTable.Update(columnsToUpdate...)
	q := gocqlx.Query(session.Query(stmt), names).BindMap(doc)
	if err := q.ExecRelease(); err != nil {
		defer InvalidateSession(ctx, ksName)
		logger.Error("Issue while trying to update the document", zap.Error(err), zap.Any("params", params), zap.Any("schema", m))
		return document.NewUpdateDocumentBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return document.NewUpdateDocumentOK().WithPayload(&models.APIResponse{Message: "Document updated"})
}
