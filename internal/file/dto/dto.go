package dto

// GenerateUploadTokenRequest generate upload token request
type GenerateUploadTokenRequest struct {
	FileName     string `json:"file_name" binding:"required" example:"photo.jpg"`
	FileSize     int64  `json:"file_size" binding:"required,gt=0" example:"1024000"`
	MimeType     string `json:"mime_type" binding:"required" example:"image/jpeg"`
	FileType     string `json:"file_type" binding:"required,oneof=image video audio file" example:"image"`
	ExpiresHours *int32 `json:"expires_hours,omitempty" example:"0"`
}

// GenerateUploadTokenResponse generate upload token response
type GenerateUploadTokenResponse struct {
	FileID    string `json:"file_id" example:"file-123"`
	UploadURL string `json:"upload_url" example:"https://minio:9000/..."`
	ExpiresIn int64  `json:"expires_in" example:"3600"`
}

// CompleteUploadRequest complete upload request
type CompleteUploadRequest struct {
	FileID string `json:"file_id" binding:"required" example:"file-123"`
}

// FileInfoResponse file info response
type FileInfoResponse struct {
	FileID        string            `json:"file_id" example:"file-123"`
	UserID        string            `json:"user_id" example:"user-123"`
	FileName      string            `json:"file_name" example:"photo.jpg"`
	FileType      string            `json:"file_type" example:"image"`
	FileSize      int64             `json:"file_size" example:"1024000"`
	MimeType      string            `json:"mime_type" example:"image/jpeg"`
	StoragePath   string            `json:"storage_path" example:"chat-file/user-123/2024-01-15/uuid.jpg"`
	ThumbnailPath string            `json:"thumbnail_path,omitempty" example:"chat-file/user-123/2024-01-15/uuid_thumb.jpg"`
	BucketName    string            `json:"bucket_name" example:"chat-file"`
	Status        int32             `json:"status" example:"1"`
	CreatedAt     int64             `json:"created_at" example:"1705315200"`
	ExpiresAt     *int64            `json:"expires_at,omitempty" example:"1705401600"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	DownloadURL   string            `json:"download_url,omitempty" example:"https://minio:9000/..."`
	ThumbnailURL  string            `json:"thumbnail_url,omitempty" example:"https://minio:9000/..."`
}

// GenerateDownloadURLRequest generate download URL request
type GenerateDownloadURLRequest struct {
	FileID         string `json:"file_id" binding:"required" example:"file-123"`
	ExpiresMinutes *int32 `json:"expires_minutes,omitempty" example:"60"`
}

// GenerateDownloadURLResponse generate download URL response
type GenerateDownloadURLResponse struct {
	DownloadURL  string `json:"download_url" example:"https://minio:9000/..."`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
	ThumbnailURL string `json:"thumbnail_url,omitempty" example:"https://minio:9000/..."`
}

// ListFilesRequest list files request
type ListFilesRequest struct {
	FileType *string `form:"file_type" example:"image"`
	Page     int     `form:"page" binding:"required,min=1" example:"1"`
	PageSize int     `form:"page_size" binding:"required,min=1,max=100" example:"20"`
}

// ListFilesResponse list files response
type ListFilesResponse struct {
	Files    []*FileInfoResponse `json:"files"`
	Total    int64               `json:"total" example:"100"`
	Page     int                 `json:"page" example:"1"`
	PageSize int                 `json:"page_size" example:"20"`
}

// DeleteFileRequest delete file request (for internal use only)
type DeleteFileRequest struct {
	FileID string `json:"file_id" binding:"required" example:"file-123"`
}

// DeleteFileResponse delete file response
type DeleteFileResponse struct {
	Success bool `json:"success" example:"true"`
}
