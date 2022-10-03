package handlers

import (
	"platform/services/storage/models"
	"platform/services/storage/restapi/operations/objectinfo"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	minio "github.com/minio/minio-go"
	zap "go.uber.org/zap"
)

// GetObjectInfoHandler implements the GET object endpoint
func GetObjectInfoHandler(params objectinfo.GetObjectInfoParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return objectinfo.NewGetObjectInfoInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// TODO: implement server side encryption
	options := minio.StatObjectOptions{}
	objName := GetIsolatedObjName(params.ObjectName, params.Owner, auth, params.Isolation)

	stats, err := client.StatObject(params.BucketName, objName, options)
	if err != nil {
		logger.Error("Object not found", zap.Error(err), zap.Any("params", params))
		return objectinfo.NewGetObjectInfoNotFound().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if !IsAuthorized(ctx, &stats, auth) {
		logger.Error("Unauthorized", zap.Error(err), zap.Any("params", params), zap.Any("objectinfo", stats))
		return objectinfo.NewGetObjectInfoUnauthorized().WithPayload(&models.APIResponse{Message: "Unauthorized"})
	}

	return objectinfo.NewGetObjectInfoOK().WithPayload(&models.ObjectInfo{
		ContentType:  stats.ContentType,
		Etag:         stats.ETag,
		LastModified: strfmt.DateTime(stats.LastModified),
		Metadata:     stats.Metadata,
		Name:         params.ObjectName,
		Owner:        &models.ObjectInfoOwner{ID: stats.Owner.ID, Name: stats.Owner.DisplayName},
		Size:         stats.Size,
		StorageClass: stats.StorageClass,
	})
}
