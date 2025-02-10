package storage

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client     *minio.Client
	bucketName string
}

func NewMinioService(endpoint, accessKey, secretKey, bucketName string) (*MinioService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists != nil || !exists {
			return nil, err
		}
	}

	return &MinioService{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *MinioService) SaveFile(id, sourceCode string, userInput string) error {
	ctx := context.Background()
	if err := s.saveFile(ctx, "source_code/"+id+".py", sourceCode); err != nil {
		return err
	}

	if err := s.saveFile(ctx, "user_input/"+id+".txt", userInput); err != nil {
		return err
	}

	return nil
}

func (s *MinioService) saveFile(ctx context.Context, objectName, content string) error {
	reader := bytes.NewReader([]byte(content))
	size := int64(len(content))

	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{})
	return err
}

func (s *MinioService) GetFile(id string) (string, string, error) {
	ctx := context.Background()
	sourceCode, err := s.getFile(ctx, "source_code/"+id+".py")
	if err != nil {
		return "", "", err
	}

	userInput, err := s.getFile(ctx, "user_input/"+id+".txt")
	if err != nil {
		return "", "", err
	}

	return sourceCode, userInput, nil
}

func (s *MinioService) getFile(ctx context.Context, objectName string) (string, error) {
	reader, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String(), nil
}
