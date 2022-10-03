package handlers

import (
	"encoding/json"
	"platform/services/storage/models"
	"platform/services/storage/restapi/operations/file"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	minio "github.com/minio/minio-go/v7"
	zap "go.uber.org/zap"
)

// AddFileHandler implements the POST object endpoint
func AddFileHandler(params file.AddFileParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return file.NewAddFileInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// TODO: unmarshal JSON directly into this object
	options := minio.PutObjectOptions{}

	// Fill the options struct
	// TODO: handle the server side encryption
	if params.UserMetadata != nil {
		options.UserMetadata = make(map[string]string)
		err = json.Unmarshal([]byte(*params.UserMetadata), &options.UserMetadata)
		if err != nil {
			logger.Error("Error reading user metadata", zap.Error(err), zap.String("meta", *params.UserMetadata))
			return file.NewAddFileInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
		}
	}
	if params.ContentType != nil {
		options.ContentType = *params.ContentType
	}
	if params.ContentEncoding != nil {
		options.ContentEncoding = *params.ContentEncoding
	}
	if params.ContentDisposition != nil {
		options.ContentDisposition = *params.ContentDisposition
	}
	if params.ContentLanguage != nil {
		options.ContentLanguage = *params.ContentLanguage
	}
	if params.CacheControl != nil {
		options.CacheControl = *params.CacheControl
	}
	if params.NumThreads != nil {
		options.NumThreads = uint(*params.NumThreads)
	}
	if params.StorageClass != nil {
		options.StorageClass = *params.StorageClass
	}
	if params.WebsiteRedirectLocation != nil {
		options.WebsiteRedirectLocation = string(*params.WebsiteRedirectLocation)
	}

	err = SetCredentials(ctx, params.Owner, &options, auth)
	if err != nil {
		return file.NewAddFileUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	objName := params.ObjectName
	if params.Isolation != nil && *params.Isolation {
		objName = options.UserMetadata["owner_name"] + "/" + params.ObjectName
		objName = strings.TrimPrefix(objName, "/")
	}

	updateInfo, err := client.PutObject(ctx, params.BucketName, objName, params.ObjectData, params.ObjectSize, options)
	if err != nil || updateInfo.Size != params.ObjectSize {
		logger.Error("Error writing object to bucket", zap.Error(err))
		return file.NewAddFileInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return file.NewAddFileCreated().WithPayload(&models.APIResponse{Message: "Object added"})
}
