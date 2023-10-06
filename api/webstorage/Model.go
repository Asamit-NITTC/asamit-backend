package webstorage

import (
	"context"
	"mime/multipart"

	"io"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type CloudStorageOriginalWebRepo struct {
	ctx    context.Context
	bucket *storage.BucketHandle
}

func InitializeCloudStorageOriginalWebRepo(c context.Context, b *storage.BucketHandle) *CloudStorageOriginalWebRepo {
	return &CloudStorageOriginalWebRepo{ctx: c, bucket: b}
}

type CloudStorageOriginalWebModel interface {
	Write(objectName string, file multipart.File) (string, error)
}

func (c CloudStorageOriginalWebRepo) Write(roomID string, file multipart.File) (string, error) {
	//roomID/UUIDにすることによってファイル名が被ることなく擬似的にディレクトリ分けしている
	fileUUID := uuid.NewString()
	objectName := roomID + "/" + fileUUID

	obj := c.bucket.Object(objectName)
	writer := obj.NewWriter(c.ctx)
	_, err := io.Copy(writer, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}
	return objectName, nil
}
