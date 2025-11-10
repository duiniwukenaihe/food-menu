package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService interface {
	GeneratePresignedUploadURL(ctx context.Context, fileName, contentType string) (uploadURL, fileURL, key string, err error)
	UploadFile(ctx context.Context, reader io.Reader, key, contentType string, size int64) (string, error)
	DeleteFile(ctx context.Context, key string) error
	GetFileURL(key string) string
}

type MinIOStorageService struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

type LocalStorageService struct {
	uploadDir string
	baseURL   string
}

// NewStorageService creates a storage service based on configuration
func NewStorageService() (StorageService, error) {
	storageType := os.Getenv("STORAGE_TYPE")
	
	if storageType == "minio" || storageType == "s3" {
		return NewMinIOStorageService()
	}
	
	return NewLocalStorageService()
}

// NewMinIOStorageService creates a new MinIO storage service
func NewMinIOStorageService() (*MinIOStorageService, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}
	
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = "media"
	}
	
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}
	
	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		
		// Set bucket policy to allow public read
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucketName)
		
		err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}
	
	return &MinIOStorageService{
		client:     minioClient,
		bucketName: bucketName,
		endpoint:   endpoint,
		useSSL:     useSSL,
	}, nil
}

// GeneratePresignedUploadURL generates a presigned URL for uploading
func (s *MinIOStorageService) GeneratePresignedUploadURL(ctx context.Context, fileName, contentType string) (uploadURL, fileURL, key string, err error) {
	ext := filepath.Ext(fileName)
	key = fmt.Sprintf("dishes/%s/%s%s", time.Now().Format("2006/01"), uuid.New().String(), ext)
	
	// Generate presigned PUT URL
	presignedURL, err := s.client.PresignedPutObject(ctx, s.bucketName, key, time.Hour)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	
	fileURL = s.GetFileURL(key)
	
	return presignedURL.String(), fileURL, key, nil
}

// UploadFile uploads a file to MinIO
func (s *MinIOStorageService) UploadFile(ctx context.Context, reader io.Reader, key, contentType string, size int64) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	
	return s.GetFileURL(key), nil
}

// DeleteFile deletes a file from MinIO
func (s *MinIOStorageService) DeleteFile(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFileURL returns the public URL for a file
func (s *MinIOStorageService) GetFileURL(key string) string {
	protocol := "http"
	if s.useSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, s.bucketName, key)
}

// NewLocalStorageService creates a local file storage service
func NewLocalStorageService() (*LocalStorageService, error) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}
	
	return &LocalStorageService{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}, nil
}

// GeneratePresignedUploadURL generates a local upload URL
func (s *LocalStorageService) GeneratePresignedUploadURL(ctx context.Context, fileName, contentType string) (uploadURL, fileURL, key string, err error) {
	ext := filepath.Ext(fileName)
	key = fmt.Sprintf("dishes/%s/%s%s", time.Now().Format("2006/01"), uuid.New().String(), ext)
	
	uploadURL = fmt.Sprintf("%s/api/v1/admin/media/upload?key=%s", s.baseURL, key)
	fileURL = s.GetFileURL(key)
	
	return uploadURL, fileURL, key, nil
}

// UploadFile uploads a file to local storage
func (s *LocalStorageService) UploadFile(ctx context.Context, reader io.Reader, key, contentType string, size int64) (string, error) {
	filePath := filepath.Join(s.uploadDir, key)
	
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	// Copy data
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	
	return s.GetFileURL(key), nil
}

// DeleteFile deletes a file from local storage
func (s *LocalStorageService) DeleteFile(ctx context.Context, key string) error {
	// Prevent directory traversal
	if strings.Contains(key, "..") {
		return fmt.Errorf("invalid file key")
	}
	
	filePath := filepath.Join(s.uploadDir, key)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}
	
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

// GetFileURL returns the public URL for a file
func (s *LocalStorageService) GetFileURL(key string) string {
	return fmt.Sprintf("%s/api/v1/media/%s", s.baseURL, key)
}
