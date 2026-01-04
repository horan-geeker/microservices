package response

import (
	"microservices/entity/model"
	"time"
)

type FileUpload struct {
	File model.File `json:"file"`
}

type UploadFile struct {
	File model.File `json:"file"`
}

type FileSummary struct {
	ID        int       `json:"id"`
	Url       string    `json:"url"`
	Type      int       `json:"type"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetFilesResp struct {
	List  []*FileSummary `json:"list"`
	Total int64          `json:"total"`
}

type FileDetailResp struct {
	File model.File `json:"file"`
}
