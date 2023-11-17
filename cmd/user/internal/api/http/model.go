package http

import (
	"github.com/gofrs/uuid"
)

// UploadFilesResponse response upload file.
type UploadFilesResponse struct {
	FileID uuid.UUID `json:"file_id"`
}
