package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"platform/services/docdb/models"
	"platform/services/docdb/restapi/operations/keyspace"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	zap "go.uber.org/zap"
)

// AddKeyspaceHandler handles the POST table request
func AddKeyspaceHandler(params keyspace.AddKeyspaceParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ks := "system"
	_, ctx, session, logger, _ := GetHandlerContext(ks, params.HTTPRequest)
	defer InvalidateSession(ctx, ks)

	ksName, err := GetKeyspaceName(*params.Keyspace.Name)
	if err != nil {
		return keyspace.NewAddKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	replication, err := json.Marshal(params.Keyspace.Replication)

	// scylla names use single quotes
	strrep := strings.ReplaceAll(string(replication), string('"'), string('\''))
	q := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = %s AND durable_writes = %t`,
		ksName, strrep, *(params.Keyspace.DurableWrites))
	err = session.Query(q).Exec()

	if err != nil {
		logger.Error("Could not create new keyspace", zap.Error(err), zap.Any("keyspace", params))
		return keyspace.NewAddKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	logger.Info("Keyspace created", zap.Any("keyspace", params))
	return keyspace.NewAddKeyspaceCreated().WithPayload(&models.APIResponse{Message: "Keyspace created"})
}

// DeleteKeyspaceHandler handles the DELETE keyspace request
func DeleteKeyspaceHandler(params keyspace.DeleteKeyspaceParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ks := "system"
	_, ctx, session, logger, _ := GetHandlerContext(ks, params.HTTPRequest)
	defer InvalidateSession(ctx, ks)

	ksName, err := GetKeyspaceName(params.KeyspaceName)
	if err != nil {
		return keyspace.NewDeleteKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// first check if the Keyspace exists
	_, err = session.KeyspaceMetadata(ksName)
	if err != nil {
		return keyspace.NewDeleteKeyspaceOK().WithPayload(&models.APIResponse{Message: "Keyspace does not exist"})
	}

	// Then delete it via CQL
	// There is no error raised via CQL commands
	q := fmt.Sprintf(`DROP KEYSPACE IF EXISTS %s`, ksName)
	err = session.Query(q).Exec()

	if err != nil {
		logger.Error("Could not delete existing keyspace", zap.Error(err), zap.String("keyspace", ksName))
		return keyspace.NewDeleteKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return keyspace.NewDeleteKeyspaceOK().WithPayload(&models.APIResponse{Message: "Keyspace deleted"})
}

// GetKeyspaceHandler handles the GET keyspace request
func GetKeyspaceHandler(params keyspace.GetKeyspaceParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ks := "system"
	_, ctx, session, logger, _ := GetHandlerContext(ks, params.HTTPRequest)

	ksName, err := GetKeyspaceName(params.KeyspaceName)
	if err != nil {
		return keyspace.NewGetKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	keyspaceMetadata, err := session.KeyspaceMetadata(ksName)
	if err != nil {
		defer InvalidateSession(ctx, ksName)
		logger.Error("Could not retrieve existing keyspace", zap.Error(err), zap.String("keyspace", ksName))
		return keyspace.NewDeleteKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	var rf int64
	if val, ok := keyspaceMetadata.StrategyOptions["replication_factor"]; ok {
		strrf := val.(string)
		rf, _ = strconv.ParseInt(strrf, 10, 64)
	}
	ret := &models.Keyspace{
		DurableWrites: &keyspaceMetadata.DurableWrites,
		Name:          &params.KeyspaceName,
		Replication: &models.Replication{
			Class:             &keyspaceMetadata.StrategyClass,
			ReplicationFactor: &rf,
		},
	}
	return keyspace.NewGetKeyspaceOK().WithPayload(ret)
}

// GetKeyspacesHandler handles the GET keyspaces request
func GetKeyspacesHandler(params keyspace.GetKeyspacesParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ks := "system"
	_, _, session, logger, _ := GetHandlerContext(ks, params.HTTPRequest)

	var foundKeyspace string
	keyspaceList := make(models.KeyspaceList, 0)
	iter := session.Query(`SELECT keyspace_name FROM system_schema.keyspaces`).Iter()
	for iter.Scan(&foundKeyspace) {
		// do not list system keyspaces
		if !strings.HasPrefix(foundKeyspace, "system") && strings.HasPrefix(foundKeyspace, os.Getenv("ENVIRONMENT")) {
			foundKeyspace = strings.TrimPrefix(foundKeyspace, os.Getenv("ENVIRONMENT")+"_")
			keyspaceList = append(keyspaceList, foundKeyspace)
		}
	}
	if err := iter.Close(); err != nil {
		logger.Fatal("Error closing namespace list iterator", zap.Error(err))
	}

	return keyspace.NewGetKeyspacesOK().WithPayload(keyspaceList)
}

// UpdateKeyspaceHandler handles the POST table request
func UpdateKeyspaceHandler(params keyspace.UpdateKeyspaceParams, auth *models.Principal) middleware.Responder {
	// Common initialization of request context
	ks := "system"
	_, ctx, session, logger, _ := GetHandlerContext(ks, params.HTTPRequest)
	defer InvalidateSession(ctx, ks)

	ksName, err := GetKeyspaceName(*params.Keyspace.Name)
	if err != nil {
		return keyspace.NewUpdateKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// first check if the Keyspace exists
	_, err = session.KeyspaceMetadata(ksName)
	if err != nil {
		logger.Error("Could not retrieve keyspace", zap.Error(err), zap.String("keyspace", ksName))
		return keyspace.NewUpdateKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	replication, err := json.Marshal(params.Keyspace.Replication)

	// scylla names use single quotes
	strrep := strings.ReplaceAll(string(replication), string('"'), string('\''))
	q := fmt.Sprintf(`ALTER KEYSPACE %s WITH replication = %s AND durable_writes = %t`,
		ksName, strrep, *(params.Keyspace.DurableWrites))
	err = session.Query(q).Exec()

	if err != nil {
		logger.Error("Could not create new keyspace", zap.Error(err), zap.Any("keyspace", params))
		return keyspace.NewUpdateKeyspaceBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return keyspace.NewUpdateKeyspaceOK().WithPayload(&models.APIResponse{Message: "Keyspace updated"})
}
