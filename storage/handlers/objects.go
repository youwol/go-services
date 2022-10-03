package handlers

import (
	"context"
	"platform/services/storage/models"
	"platform/services/storage/restapi/operations/objects"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/minio/minio-go"
	zap "go.uber.org/zap"
)

// GetObjectsHandler implements the GET objects endpoint
func GetObjectsHandler(params objects.GetObjectsParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return objects.NewGetObjectsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	metaChannel := make(chan models.ObjectInfo, 10)
	go func() {
		defer close(metaChannel)

		// Create a done channel to control 'ListObjects' go routine.
		doneChannel := make(chan struct{})
		// Indicate to our routine to exit cleanly upon return.
		defer close(doneChannel)

		var prefix = string("")
		if params.Prefix != nil {
			prefix = *params.Prefix
		}

		var recursive = true
		if params.Recursive != nil {
			recursive = *params.Recursive
		}

		// If we have set a recursive group search with prefix, then we cannot handle all subgroups with only one prefix
		// So we search without prefix, and then filter the results
		var resultsFilter bool = false
		search := prefix
		search = GetIsolatedObjName(prefix, params.Owner, auth, params.Isolation)
		if recursive && params.Owner != nil && len(*params.Owner) > 0 {
			resultsFilter = true
		}

		for obj := range client.ListObjectsV2(params.BucketName, search, recursive, doneChannel) {
			if obj.Err != nil {
				logger.Error("Error retrieving one of the objects", zap.Error(obj.Err), zap.Any("object", obj))
				continue
			}

			if resultsFilter {
				// check the file does not belong to a sub-group
				cnt := false
				for _, grp := range auth.MemberOf {
					if grp.Path == *params.Owner {
						continue
					}

					grpWithoutSlash := grp.Path[1:]
					if !strings.HasPrefix(obj.Key, *params.Owner) && strings.HasPrefix(obj.Key, grpWithoutSlash) {
						cnt = true
					}
				}
				if cnt {
					continue
				}
			}
			// Currently the ListObjectsV2 returns empty metadata... Minio enh on it way
			// Meanwhile, we need to fetch the meta explicitely on each object :'(
			// PERF WARNING
			meta, err := client.StatObject(params.BucketName, obj.Key, minio.StatObjectOptions{})
			if err == nil && IsAuthorized(ctx, &meta, auth) {
				metaChannel <- models.ObjectInfo{
					Etag:         meta.ETag,
					LastModified: strfmt.DateTime(meta.LastModified),
					Metadata:     meta.Metadata,
					Name:         meta.Key,
					Owner:        &models.ObjectInfoOwner{ID: meta.Owner.ID, Name: meta.Owner.DisplayName},
					Size:         meta.Size,
					StorageClass: meta.StorageClass,
					ContentType:  meta.ContentType,
				}
			} else {
				logger.Warn("One of the objects was removed from the list (unauthorized)", zap.String("name", obj.Key), zap.Any("object", obj.Metadata), zap.Error(err))
			}
		}
	}()

	ret := make([]*models.ObjectInfo, 0)
	for r := range metaChannel {
		entry := r
		ret = append(ret, &entry)
	}
	if len(ret) == 0 {
		logger.Warn("Objects not found", zap.Any("params", params))
		// put the error as payload, if any
		return objects.NewGetObjectsNotFound()
	}

	return objects.NewGetObjectsOK().WithPayload(ret)
}

// DeleteObjects removes objects from a bucket
func DeleteObjects(client *minio.Client, logger zap.Logger, params objects.DeleteObjectsParams, auth *models.Principal) error {
	var prefix = string("")
	if params.Prefix != nil {
		prefix = *params.Prefix
	}

	var recursive = true
	if params.Recursive != nil {
		recursive = *params.Recursive
	}

	// If we have set a recursive group search with prefix, then we cannot handle all subgroups with only one prefix
	// So we search without prefix, and then filter the results
	var resultsFilter bool = false
	search := GetIsolatedObjName(prefix, params.Owner, auth, params.Isolation)
	if recursive && params.Owner != nil && len(*params.Owner) > 0 {
		resultsFilter = true
	}

	// Create a deleteChannel to store objects to delete
	deleteChannel := make(chan string)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(deleteChannel)

		doneCh := make(chan struct{})

		// Indicate to our routine to exit cleanly upon return.
		defer close(doneCh)

		// List all objects from a bucket-name with a matching prefix.
		for object := range client.ListObjects(params.BucketName, search, recursive, doneCh) {
			if object.Err != nil {
				logger.Error("Error retrieving one of the objects", zap.Error(object.Err), zap.Any("object", object))
			}

			if resultsFilter {
				// check the file does not belong to a sub-group
				cnt := false
				for _, grp := range auth.MemberOf {
					if grp.Path == *params.Owner {
						continue
					}

					grpWithoutSlash := grp.Path[1:]
					if !strings.HasPrefix(object.Key, *params.Owner) && strings.HasPrefix(object.Key, grpWithoutSlash) {
						cnt = true
					}
				}
				if cnt {
					continue
				}
			}

			meta, err := client.StatObject(params.BucketName, object.Key, minio.StatObjectOptions{})
			if err == nil && IsAuthorized(context.Background(), &meta, auth) {
				deleteChannel <- object.Key
			} else {
				logger.Error("Error deleting one of the objects (unauthorized)", zap.Error(err), zap.Any("object", meta))
			}
		}
	}()

	errorCh := client.RemoveObjects(params.BucketName, deleteChannel)
	var err error
	for e := range errorCh {
		logger.Error("Error deleting one of the objects", zap.Error(e.Err), zap.Any("object", e))
		err = e.Err
	}
	return err
}

// DeleteObjectsHandler implements the DELETE objects endpoint
func DeleteObjectsHandler(params objects.DeleteObjectsParams, auth *models.Principal) middleware.Responder {
	_, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return objects.NewDeleteObjectsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	err = DeleteObjects(client, logger, params, auth)

	if err != nil {
		return objects.NewDeleteObjectsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return objects.NewDeleteObjectsOK().WithPayload(&models.APIResponse{Message: "Objects deleted"})
}
