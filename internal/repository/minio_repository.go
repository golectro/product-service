package repository

import (
	"bytes"
	"context"
	"time"

	"golectro-product/internal/model"

	"github.com/minio/minio-go/v7"
)

type MinioRepository struct {
	Client *minio.Client
}

func NewMinioRepository(client *minio.Client) *MinioRepository {
	return &MinioRepository{Client: client}
}

func (r *MinioRepository) UploadFile(ctx context.Context, input model.UploadFileInput) error {
	_, err := r.Client.PutObject(ctx, input.Bucket, input.ObjectKey, bytes.NewReader(input.Content), int64(len(input.Content)), minio.PutObjectOptions{
		ContentType: input.ContentType,
	})
	return err
}

func (r *MinioRepository) GeneratePresignedURL(ctx context.Context, input model.PresignedURLInput) (string, error) {
	reqParams := make(map[string][]string)
	presignedURL, err := r.Client.PresignedGetObject(ctx, input.Bucket, input.ObjectKey, time.Duration(input.Expiry)*time.Second, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (r *MinioRepository) DeleteFile(ctx context.Context, bucket, objectKey string) error {
	return r.Client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{})
}

func (r *MinioRepository) GetObject(ctx context.Context, bucket, objectKey string) (*minio.Object, error) {
	object, err := r.Client.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
