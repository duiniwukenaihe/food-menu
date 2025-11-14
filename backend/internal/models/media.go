package models

import "time"

// Media represents uploaded media metadata in the database
type Media struct {
    ID           uint       `json:"id" gorm:"primaryKey"`
    Key          string     `json:"key" gorm:"type:varchar(500);uniqueIndex;not null"`
    URL          string     `json:"url" gorm:"type:varchar(500);not null"`
    FileName     string     `json:"fileName" gorm:"type:varchar(255);not null"`
    ContentType  string     `json:"contentType" gorm:"type:varchar(100);not null"`
    Size         int64      `json:"size" gorm:"not null"`
    EntityType   string     `json:"entityType" gorm:"type:varchar(50);index"` // e.g., "dish", "user"
    EntityID     *uint      `json:"entityId" gorm:"index"`
    IsPrimary    bool       `json:"isPrimary" gorm:"default:false"`
    AltText      string     `json:"altText" gorm:"type:varchar(255)"`
    CreatedBy    *uint      `json:"createdBy" gorm:"index"`
    UpdatedBy    *uint      `json:"updatedBy" gorm:"index"`
    CreatedAt    time.Time  `json:"createdAt"`
    UpdatedAt    time.Time  `json:"updatedAt"`
}


// UploadURLRequest represents a request for a presigned upload URL
type UploadURLRequest struct {
    FileName    string `json:"fileName" binding:"required,max=255"`
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

// DeleteMediaRequest represents a request to delete media
type DeleteMediaRequest struct {
    Key string `json:"key" binding:"required"`
}
