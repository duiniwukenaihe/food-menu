package models

// UploadURLRequest represents a request for a presigned upload URL
type UploadURLRequest struct {
	FileName    string `json:"fileName" binding:"required"`
	ContentType string `json:"contentType" binding:"required"`
}

// UploadURLResponse represents the response with presigned upload URL
type UploadURLResponse struct {
	UploadURL string `json:"uploadUrl"`
	FileURL   string `json:"fileUrl"`
	Key       string `json:"key"`
}

// MediaFile represents uploaded media information
type MediaFile struct {
	URL         string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}
