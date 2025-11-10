//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"example.com/app/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestMinIOStorageService(t *testing.T) {
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		t.Skip("MINIO_ENDPOINT not set, skipping MinIO tests")
	}

	os.Setenv("STORAGE_TYPE", "minio")

	storageService, err := services.NewStorageService()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("Generate Presigned Upload URL", func(t *testing.T) {
		uploadURL, fileURL, key, err := storageService.GeneratePresignedUploadURL(ctx, "test.jpg", "image/jpeg")

		assert.NoError(t, err)
		assert.NotEmpty(t, uploadURL)
		assert.NotEmpty(t, fileURL)
		assert.NotEmpty(t, key)
		assert.Contains(t, key, "dishes/")
		assert.Contains(t, key, ".jpg")
	})

	t.Run("Upload and Retrieve File", func(t *testing.T) {
		content := []byte("test file content for MinIO")
		reader := bytes.NewReader(content)

		fileURL, err := storageService.UploadFile(ctx, reader, "dishes/2024/01/test-minio.jpg", "image/jpeg", int64(len(content)))

		assert.NoError(t, err)
		assert.NotEmpty(t, fileURL)
		assert.Contains(t, fileURL, "dishes/2024/01/test-minio.jpg")
	})

	t.Run("Delete File", func(t *testing.T) {
		content := []byte("test file to delete")
		reader := bytes.NewReader(content)
		key := "dishes/2024/01/delete-test.jpg"

		_, err := storageService.UploadFile(ctx, reader, key, "image/jpeg", int64(len(content)))
		assert.NoError(t, err)

		err = storageService.DeleteFile(ctx, key)
		assert.NoError(t, err)
	})
}

func TestLocalStorageService(t *testing.T) {
	os.Setenv("STORAGE_TYPE", "local")
	os.Setenv("UPLOAD_DIR", "./test_uploads_local")

	defer os.RemoveAll("./test_uploads_local")

	storageService, err := services.NewStorageService()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("Generate Local Upload URL", func(t *testing.T) {
		uploadURL, fileURL, key, err := storageService.GeneratePresignedUploadURL(ctx, "test.jpg", "image/jpeg")

		assert.NoError(t, err)
		assert.NotEmpty(t, uploadURL)
		assert.NotEmpty(t, fileURL)
		assert.NotEmpty(t, key)
		assert.Contains(t, uploadURL, "/admin/media/upload")
	})

	t.Run("Upload File Locally", func(t *testing.T) {
		content := []byte("test local file content")
		reader := bytes.NewReader(content)

		fileURL, err := storageService.UploadFile(ctx, reader, "dishes/2024/01/test-local.jpg", "image/jpeg", int64(len(content)))

		assert.NoError(t, err)
		assert.NotEmpty(t, fileURL)
		assert.Contains(t, fileURL, "dishes/2024/01/test-local.jpg")

		_, err = os.Stat("./test_uploads_local/dishes/2024/01/test-local.jpg")
		assert.NoError(t, err)
	})

	t.Run("Delete Local File", func(t *testing.T) {
		content := []byte("test file to delete locally")
		reader := bytes.NewReader(content)
		key := "dishes/2024/01/delete-test-local.jpg"

		_, err := storageService.UploadFile(ctx, reader, key, "image/jpeg", int64(len(content)))
		assert.NoError(t, err)

		err = storageService.DeleteFile(ctx, key)
		assert.NoError(t, err)

		_, err = os.Stat("./test_uploads_local/" + key)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Prevent Directory Traversal", func(t *testing.T) {
		content := []byte("malicious content")
		reader := bytes.NewReader(content)
		key := "../../../etc/passwd"

		err := storageService.DeleteFile(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "invalid")
	})
}

func TestStorageServiceVideoSupport(t *testing.T) {
	os.Setenv("STORAGE_TYPE", "local")
	os.Setenv("UPLOAD_DIR", "./test_uploads_video")

	defer os.RemoveAll("./test_uploads_video")

	storageService, err := services.NewStorageService()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("Upload Video File", func(t *testing.T) {
		content := []byte("fake video content")
		reader := bytes.NewReader(content)

		fileURL, err := storageService.UploadFile(ctx, reader, "dishes/2024/01/test-video.mp4", "video/mp4", int64(len(content)))

		assert.NoError(t, err)
		assert.NotEmpty(t, fileURL)
		assert.Contains(t, fileURL, "test-video.mp4")
	})
}
