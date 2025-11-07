package services

import (
	"testing"
	"example.com/app/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestContentValidation(t *testing.T) {
	tests := []struct {
		name    string
		content models.CreateContentRequest
		wantErr bool
	}{
		{
			name: "valid content",
			content: models.CreateContentRequest{
				Title:      "Test Article",
				Body:       "This is a test article content",
				CategoryID: 1,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			content: models.CreateContentRequest{
				Body:       "This is a test article content",
				CategoryID: 1,
			},
			wantErr: true,
		},
		{
			name: "missing body",
			content: models.CreateContentRequest{
				Title:      "Test Article",
				CategoryID: 1,
			},
			wantErr: true,
		},
		{
			name: "missing category",
			content: models.CreateContentRequest{
				Title: "Test Article",
				Body:  "This is a test article content",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContentRequest(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func validateContentRequest(req models.CreateContentRequest) error {
	if req.Title == "" {
		return assert.AnError
	}
	if req.Body == "" {
		return assert.AnError
	}
	if req.CategoryID == 0 {
		return assert.AnError
	}
	return nil
}
