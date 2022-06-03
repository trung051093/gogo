package storage

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService interface {
	FGetObject(ctx context.Context, bucketName string, objectName string, filePath string, options minio.GetObjectOptions) error
	FPutObject(ctx context.Context, bucketName string, objectName string, fileName string, option minio.PutObjectOptions) (info minio.UploadInfo, err error)
	PresignedPutObject(ctx context.Context, bucketName string, objectName string, expires time.Duration) (*url.URL, error)
	PresignedGetObject(ctx context.Context, bucketName string, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error)
}

type StorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type storageService struct {
	client *minio.Client
}

func NewStorage(config *StorageConfig) (*storageService, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		log.Println("Cannot connect minio service: ", err)
		return nil, err
	}

	return &storageService{client: client}, nil
}

func (s *storageService) CreateBucket(ctx context.Context, bucketName string, location string) error {
	err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s.client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	return nil
}

func (s *storageService) FGetObject(ctx context.Context, bucketName string, objectName string, filePath string, options minio.GetObjectOptions) error {
	return s.client.FGetObject(ctx, bucketName, objectName, filePath, options)
}

func (s *storageService) FPutObject(ctx context.Context, bucketName string, objectName string, fileName string, options minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	return s.client.FPutObject(ctx, bucketName, objectName, fileName, options)
}

func (s *storageService) PresignedPutObject(ctx context.Context, bucketName string, objectName string, expires time.Duration) (*url.URL, error) {
	return s.client.PresignedPutObject(ctx, bucketName, objectName, expires)
}

func (s *storageService) PresignedGetObject(ctx context.Context, bucketName string, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	return s.client.PresignedGetObject(ctx, bucketName, objectName, expires, reqParams)
}
