package webstorage

import (
	"context"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

type CloudStorageWebRepo struct {
	ctx    context.Context
	bucket *storage.BucketHandle
}

func InitializeCloudStorageRepo(c context.Context, b *storage.BucketHandle) *CloudStorageWebRepo {
	return &CloudStorageWebRepo{ctx: c, bucket: b}
}

type CloudStorageWebModel interface {
	Write(objectName string, file multipart.File) (string, error)
}

func (c CloudStorageWebRepo) Write(roomID string, file multipart.File) (string, error) {
	return "", nil
}
