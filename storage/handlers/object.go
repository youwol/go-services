package handlers

import (
	"encoding/base64"
	"io/ioutil"
	"platform/services/storage/models"
	"platform/services/storage/restapi/operations/object"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	minio "github.com/minio/minio-go/v7"
	zap "go.uber.org/zap"
)

// AddObjectHandler implements the POST object endpoint
func AddObjectHandler(params object.AddObjectParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return object.NewAddObjectInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// TODO: unmarshal JSON directly into this object
	options := minio.PutObjectOptions{}
	options.UserMetadata = make(map[string]string)
	// options.ContentEncoding = "utf-8"

	// Fill the options struct
	if params.Data.Options != nil {
		if params.Data.Options.UserMetadata != nil {
			options.UserMetadata = params.Data.Options.UserMetadata.(map[string]string)
		}
		if len(params.Data.Options.ContentType) > 0 {
			options.ContentType = params.Data.Options.ContentType
		}
		if len(params.Data.Options.ContentEncoding) > 0 {
			options.ContentEncoding = params.Data.Options.ContentEncoding
		}
		if len(params.Data.Options.ContentDisposition) > 0 {
			options.ContentDisposition = params.Data.Options.ContentDisposition
		}
		if len(params.Data.Options.ContentLanguage) > 0 {
			options.ContentLanguage = params.Data.Options.ContentLanguage
		}
		if len(params.Data.Options.CacheControl) > 0 {
			options.CacheControl = params.Data.Options.CacheControl
		}
		if params.Data.Options.NumThreads > 0 {
			options.NumThreads = uint(params.Data.Options.NumThreads)
		}
		if params.Data.Options.StorageClass != nil {
			options.StorageClass = *params.Data.Options.StorageClass
		}
		if len(params.Data.Options.WebsiteRedirectLocation) > 0 {
			options.WebsiteRedirectLocation = string(params.Data.Options.WebsiteRedirectLocation)
		}
	}

	// Set the ownership metadata, unfortunately the minio custom metadata are flat (no nested structs)
	// So we have to flattenize the list of groups
	var owner *string
	if params.Owner != nil {
		owner = params.Owner
	}
	err = SetCredentials(ctx, owner, &options, auth)
	if err != nil {
		return object.NewAddObjectUnauthorized().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	objName := GetIsolatedObjName(params.Data.Object.Name, owner, auth, params.Isolation)

	// If no isolation, we must make sure that we do not overwrite someone else's data
	if params.Isolation != nil && *params.Isolation == false {
		obj, err := client.StatObject(ctx, params.BucketName, objName, minio.StatObjectOptions{})
		if err == nil && !IsAuthorized(ctx, &obj, auth) {
			return object.NewAddObjectUnauthorized().WithPayload(&models.APIResponse{Message: "The object already exists and does not belong to you"})
		}
	}

	rd := base64.NewDecoder(base64.StdEncoding, strings.NewReader(params.Data.Object.Data.String()))
	_, err = client.PutObject(ctx, params.BucketName, objName, rd, params.Data.Object.Size, options)
	if err != nil {
		logger.Error("Error writing object to bucket", zap.Error(err))
		return object.NewAddObjectInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return object.NewAddObjectCreated().WithPayload(&models.APIResponse{Message: "Object added"})
}

// GetObjectHandler implements the GET object endpoint
func GetObjectHandler(params object.GetObjectParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return object.NewGetObjectInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// TODO: implement server side encryption
	options := minio.GetObjectOptions{}

	objName := GetIsolatedObjName(params.ObjectName, params.Owner, auth, params.Isolation)

	obj, err := client.GetObject(ctx, params.BucketName, objName, options)
	if err != nil {
		logger.Error("Could not retrieve the stored object", zap.Error(err), zap.Any("params", params))
		return object.NewGetObjectInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	stats, err := obj.Stat()
	if stats.Size == 0 || err != nil {
		logger.Error("Object not found", zap.Error(err), zap.Any("params", params))
		return object.NewGetObjectNotFound().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if !IsAuthorized(ctx, &stats, auth) {
		logger.Error("Unauthorized", zap.Error(err), zap.Any("params", params))
		return object.NewGetObjectUnauthorized().WithPayload(&models.APIResponse{Message: "Unauthorized"})
	}

	// bof bof - danger with huge data - TODO : find a way to stream that with ioReadCloser
	var b strfmt.Base64
	b, err = ioutil.ReadAll(obj)
	if err != nil {
		logger.Error("Could not read the raw buffer", zap.Error(err))
		return object.NewGetObjectInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	b, err = b.MarshalJSON()
	if err != nil {
		logger.Error("Could not b64 encode the response", zap.Error(err))
		return object.NewGetObjectInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer obj.Close()
	return object.NewGetObjectOK().WithPayload(b)
}

// DeleteObjectHandler implements the DELETE object endpoint
func DeleteObjectHandler(params object.DeleteObjectParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return object.NewDeleteObjectInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	// Authorize deletion
	options := minio.StatObjectOptions{}
	objName := GetIsolatedObjName(params.ObjectName, params.Owner, auth, params.Isolation)
	stats, err := client.StatObject(ctx, params.BucketName, objName, options)
	if err != nil {
		logger.Error("Could not delete the stored object", zap.Error(err), zap.Any("params", params))
		return object.NewDeleteObjectNotFound().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	if !IsAuthorized(ctx, &stats, auth) {
		logger.Error("Unauthorized", zap.Any("params", params), zap.Any("objectinfo", stats))
		return object.NewGetObjectUnauthorized().WithPayload(&models.APIResponse{Message: "Unauthorized"})
	}

	err = client.RemoveObject(ctx, params.BucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		logger.Error("Could not delete the stored object", zap.Error(err), zap.Any("params", params))
		return object.NewDeleteObjectNotFound().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return object.NewDeleteObjectOK().WithPayload(&models.APIResponse{Message: "Object deleted"})
}
