package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"platform/services/docdb/models"
	"platform/services/docdb/restapi/operations/table"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"gitlab.com/youwol/platform/libs/go-libs/utils"
	zap "go.uber.org/zap"
)

// TODO: the gocqlx api now has a Table object

// super nasty function to transform the JSON table description into a CQL query
func queryFromTableOptions(ctx context.Context, tableOptions *models.TableOptions, clusteringColumns *[]string) (string, error) {
	var q string
	var err error
	if tableOptions != nil {
		q += " WITH "
		bDone := false
		if tableOptions.ClusteringOrder != nil && clusteringColumns != nil {
			q += "CLUSTERING ORDER BY ("
			for _, o := range tableOptions.ClusteringOrder {
				if !contains(*clusteringColumns, *o.Name) {
					err = errors.New("Incorrect clustering order. Column must be a clustering column")
					logger := utils.ContextLogger(ctx)
					logger.Error("Clustering configuration", zap.Error(err), zap.String("col", *o.Name))
				}
				q += *o.Name + " " + *o.Order + ","
			}
			q = strings.TrimRight(q, ",")
			q += ")"
			bDone = true
			tableOptions.ClusteringOrder = tableOptions.ClusteringOrder[:0]
		}

		if tableOptions.Compaction != nil {
			if bDone {
				q += " AND "
			}

			compaction, _ := json.Marshal(tableOptions.Compaction)
			strrep := strings.ReplaceAll(string(compaction), string('"'), string('\''))
			q += "compaction = " + strrep
			bDone = true
			tableOptions.Compaction = nil
		}

		if tableOptions.Compression != nil {
			if bDone {
				q += " AND "
			}

			compression, _ := json.Marshal(tableOptions.Compression)
			strrep := strings.ReplaceAll(string(compression), string('"'), string('\''))
			q += "compression = " + strrep
			bDone = true
			tableOptions.Compression = nil
		}

		if len(tableOptions.Comment) > 0 {
			if bDone {
				q += " AND "
			}
			q += "comment = '" + tableOptions.Comment + "'"
			bDone = true
			tableOptions.Comment = ""
		}

		if tableOptions.SpeculativeRetry != nil {
			if bDone {
				q += " AND "
			}
			q += "speculative_retry = '" + *tableOptions.SpeculativeRetry + "'"
			bDone = true
			tableOptions.SpeculativeRetry = nil
		}

		// Now that we have extracted nested structs, let's transform the flat options
		options, _ := json.Marshal(tableOptions)
		stropt := strings.ReplaceAll(string(options), "\"", "")
		stropt = strings.ReplaceAll(stropt, ",", " AND ")
		stropt = strings.ReplaceAll(stropt, ":", " = ")
		stropt = strings.TrimPrefix(stropt, "{")
		stropt = strings.TrimRight(stropt, "}")
		if bDone && len(stropt) > 0 {
			q += " AND "
		}
		q += stropt
	}

	return q, err
}

// AddTableHandler handles the POST table request
func AddTableHandler(params table.AddTableParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)
	defer InvalidateSession(ctx, ksName)

	if err != nil {
		return table.NewAddTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Add ownership columns as partition key
	ownerColType := "text"
	params.Table.Columns = append(params.Table.Columns,
		&models.Column{Name: &OwnerIDColumn, Type: &ownerColType},
		&models.Column{Name: &OwnerNameColumn, Type: &ownerColType},
		&models.Column{Name: &OwnerKindColumn, Type: &ownerColType},
	)
	params.Table.PartitionKey = append(params.Table.PartitionKey, OwnerIDColumn, OwnerNameColumn, OwnerKindColumn)

	var columns string
	for _, col := range params.Table.Columns {
		colProp := ""
		if col.Static != nil { // no ternary operator in golang
			if *col.Static {
				colProp = " STATIC"
			}
		}
		if col.PrimaryKey != nil { // no ternary operator in golang
			if *col.PrimaryKey {
				colProp += " PRIMARY KEY"
			}
		}
		column := fmt.Sprintf("%s %s%s, ", *col.Name, *col.Type, colProp)
		columns += column
	}

	// Build the primary key, made of one ore more partition keys
	var primaryKey string
	if len(params.Table.PartitionKey) == 1 {
		primaryKey = params.Table.PartitionKey[0]
	} else {
		primaryKey = "("
		for _, col := range params.Table.PartitionKey {
			primaryKey += col + ", "
		}
		// replace last ", " with the closing parenthesis
		primaryKey = strings.TrimRight(primaryKey, ", ")
		primaryKey += ")"
	}

	// Append optional clustering keys
	for _, col := range params.Table.ClusteringColumns {
		primaryKey += ", " + col
	}

	// Let's build the query
	q := "CREATE TABLE IF NOT EXISTS"
	q = fmt.Sprintf(`%s %s (%sPRIMARY KEY (%s))`,
		q,
		*params.Table.Name,
		columns,
		primaryKey)

	opts, err := queryFromTableOptions(ctx, params.Table.TableOptions, &params.Table.ClusteringColumns)
	if err != nil {
		return table.NewAddTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	q += opts

	logger.Debug("built table query", zap.String("query", q))
	err = session.Query(q).Exec()
	if err != nil {
		logger.Error("Could not create new table", zap.Error(err), zap.Any("params", params), zap.String("query", q))
		return table.NewAddTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return table.NewAddTableCreated().WithPayload(&models.APIResponse{Message: "Table created"})
}

// UpdateTableHandler handles the PUT table request
func UpdateTableHandler(params table.UpdateTableParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)
	defer InvalidateSession(ctx, ksName)

	if err != nil {
		return table.NewUpdateTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	q := "ALTER TABLE "
	q += *params.TableUpdate.Name

	var addColumns string
	for _, col := range params.TableUpdate.AddColumns {
		// Ignoring STATIC and PRIMARY KEY (not available for ALTER)
		column := fmt.Sprintf("%s %s, ", *col.Name, *col.Type)
		addColumns += column
	}
	if len(addColumns) > 0 {
		q += " ADD (" + addColumns
		q = strings.TrimRight(q, ", ")
		q += ")"
	}

	var dropColumns string
	for _, col := range params.TableUpdate.DropColumns {
		// Ignoring STATIC and PRIMARY KEY (not available for ALTER)
		dropColumns += col + ", "
	}
	if len(dropColumns) > 0 {
		q += " DROP (" + dropColumns
		q = strings.TrimRight(q, ", ")
		q += ")"
	}

	opts, err := queryFromTableOptions(ctx, params.TableUpdate.TableOptions, nil)
	if err != nil {
		return table.NewAddTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	q += opts

	logger.Debug("built table query", zap.String("query", q))
	err = session.Query(q).Exec()
	if err != nil {
		logger.Error("Could not update table", zap.Error(err), zap.Any("params", params))
		return table.NewUpdateTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return table.NewUpdateTableOK().WithPayload(&models.APIResponse{Message: "Table updated"})
}

// DeleteTableHandler handles the DELETE table request
func DeleteTableHandler(params table.DeleteTableParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, ctx, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)
	if err != nil {
		return table.NewDeleteTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	defer InvalidateSession(ctx, ksName)

	// first check if the Keyspace exists
	keyspaceMetadata, err := session.KeyspaceMetadata(ksName)
	_, found := keyspaceMetadata.Tables[params.TableName]
	if found {
		q := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, params.TableName)
		err = session.Query(q).Exec()

		if err != nil {
			logger.Error("Could not delete existing table", zap.Error(err), zap.Any("params", params))
			return table.NewDeleteTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
		}
		return table.NewDeleteTableOK().WithPayload(&models.APIResponse{Message: "Table deleted"})
	}

	return table.NewDeleteTableOK().WithPayload(&models.APIResponse{Message: "Table does not exist"})
}

// GetTableHandler handles the GET table request
func GetTableHandler(params table.GetTableParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, _, session, logger, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return table.NewGetTableBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	keyspaceMetadata, _ := session.KeyspaceMetadata(ksName)
	tmd, found := keyspaceMetadata.Tables[params.TableName]
	if !found {
		return table.NewGetTableBadRequest().WithPayload(&models.APIResponse{Message: "Table not found"})
	}

	logger.Debug("retrieved table data", zap.Any("table", tmd))

	// Unfortunately we need to map manually from CQL to our API
	// GoCQL does not have the same exact structure as the Cassandra API
	columns := make([]*models.Column, len(tmd.Columns)-3) // owner id, name, kind
	partitionKey := make([]string, 0)
	clusteringCols := make([]string, 0)
	i := 0
	for _, val := range tmd.Columns {
		// Do not describe internal columns
		if val.Name == OwnerIDColumn || val.Name == OwnerNameColumn || val.Name == OwnerKindColumn {
			continue
		}
		columns[i] = &models.Column{Name: &val.Name, Type: &val.Type}
		if val.Kind == 1 { // reversed engineered the magic number
			partitionKey = append(partitionKey, val.Name)
		} else if val.Kind == 2 {
			clusteringCols = append(clusteringCols, val.Name)
		}
		i++
	}

	optionsJSON, err := json.Marshal(&tmd.Options)
	if err != nil {
		return table.NewGetTableInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	optionsModel := models.TableOptions{}
	err = json.Unmarshal(optionsJSON, &optionsModel)
	if err != nil {
		return table.NewGetTableInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	t := models.Table{
		Name:              &tmd.Name,
		Columns:           columns,
		ClusteringColumns: clusteringCols,
		PartitionKey:      partitionKey,
		TableOptions:      &optionsModel,
	}

	return table.NewGetTableOK().WithPayload(&t)
}

// GetTablesHandler handles the GET tables request
func GetTablesHandler(params table.GetTablesParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ksName, _, session, _, err := GetHandlerContext(params.KeyspaceName, params.HTTPRequest)

	if err != nil {
		return table.NewGetTablesBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	keyspaceMetadata, _ := session.KeyspaceMetadata(ksName)
	tableList := make(models.TableList, 0)
	for _, t := range keyspaceMetadata.Tables {
		tableList = append(tableList, t.Name)
	}

	return table.NewGetTablesOK().WithPayload(tableList)
}
