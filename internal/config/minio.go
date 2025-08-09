package config

import (
	"context"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewMinioClient(viper *viper.Viper, logger *logrus.Logger) *minio.Client {
	endpoint := viper.GetString("MINIO_ENDPOINT")
	accessKey := viper.GetString("MINIO_ACCESS_KEY")
	secretKey := viper.GetString("MINIO_SECRET_KEY")
	useSSL := viper.GetBool("MINIO_USE_SSL")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	if _, err := client.ListBuckets(context.Background()); err != nil {
		logger.Fatalf("Failed to connect to MinIO server: %v", err)
	}

	logger.Info("MinIO connected successfully")

	bucketsStr := viper.GetString("MINIO_BUCKETS")
	if bucketsStr == "" {
		logger.Fatal("MINIO_BUCKETS is required in environment variables")
	}

	buckets := strings.Split(bucketsStr, ",")
	for _, bucket := range buckets {
		bucket = strings.TrimSpace(bucket)
		if bucket == "" {
			continue
		}

		exists, err := client.BucketExists(context.Background(), bucket)
		if err != nil {
			logger.Fatalf("Failed to check if bucket '%s' exists: %v", bucket, err)
		}

		if !exists {
			err := client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
			if err != nil {
				logger.Fatalf("Failed to create bucket '%s': %v", bucket, err)
			}
			logger.Infof("Created new MinIO bucket: %s", bucket)
		} else {
			logger.Infof("MinIO bucket already exists: %s", bucket)
		}
	}

	return client
}
