package webstorage

import (
	"context"
	"mime/multipart"

	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"

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
	//roomID/UUIDにすることによってファイル名が被ることなく擬似的にディレクトリ分けしている
	fileUUID := uuid.NewString()
	objectName := roomID + "/" + fileUUID

	ctx := c.ctx
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := c.bucket.Object(objectName)

	o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return objectName, nil
}
