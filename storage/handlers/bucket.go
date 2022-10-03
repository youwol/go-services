package handlers

import (
	"github.com/minio/minio-go/v7"
	"platform/services/storage/models"
	"platform/services/storage/restapi/operations/bucket"
	"platform/services/storage/restapi/operations/objects"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	zap "go.uber.org/zap"
)

// AddBucketHandler handles the POST bucket request
func AddBucketHandler(params bucket.AddBucketParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	// defer InvalidateClient()
	if err != nil {
		return bucket.NewAddBucketInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	makeBucketOptions := minio.MakeBucketOptions{Region: region, ObjectLocking: true}

	err = client.MakeBucket(ctx, *(params.Bucket.Name), makeBucketOptions)
	if err != nil {
		logger.Error("Could not create bucket", zap.Error(err))
		return bucket.NewAddBucketBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return bucket.NewAddBucketCreated().WithPayload(&models.APIResponse{Message: "Bucket created"})
}

// DeleteBucketHandler handles the DELETE bucket request
func DeleteBucketHandler(params bucket.DeleteBucketParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	// defer InvalidateClient()
	if err != nil {
		return bucket.NewDeleteBucketInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	bExists, err := client.BucketExists(ctx, params.BucketName)
	if !bExists {
		logger.Error("Could not find bucket")
		return bucket.NewDeleteBucketNotFound()
	}

	// Remove bucket contents if forceNotEmpty is passed
	if params.ForceNotEmpty != nil && *params.ForceNotEmpty == true {
		deleteParams := &objects.DeleteObjectsParams{BucketName: params.BucketName}
		DeleteObjects(ctx, client, logger, *deleteParams, auth)
	}

	err = client.RemoveBucket(ctx, params.BucketName)
	if err != nil {
		logger.Error("Could not delete bucket", zap.Error(err))
		return bucket.NewDeleteBucketBadRequest().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return bucket.NewDeleteBucketOK().WithPayload(&models.APIResponse{Message: "Bucket deleted"})
}

// GetBucketsHandler handles the GET buckets request
func GetBucketsHandler(params bucket.GetBucketsParams, auth *models.Principal) middleware.Responder {
	ctx, client, logger, err := GetHandlerContext(params.HTTPRequest)
	if err != nil {
		return bucket.NewGetBucketsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		logger.Error("Could not list buckets", zap.Error(err))
		return bucket.NewGetBucketsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ret := make([]*models.Bucket, len(buckets))
	for i, b := range buckets {
		// bug if using b.Name...don't know why
		ret[i] = &models.Bucket{Name: &buckets[i].Name, CreationDate: strfmt.DateTime(b.CreationDate)}
	}

	return bucket.NewGetBucketsOK().WithPayload(ret)
}
