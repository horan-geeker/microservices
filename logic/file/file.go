package file

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"microservices/cache"
	entity "microservices/entity/model"
	"microservices/entity/response"
	"microservices/model"
	"microservices/service"
	"net/http"
	"strings"
	"time"
)

type Logic interface {
	UploadFile(ctx context.Context, userId int, fileContent []byte, fileType int, fileOriginName string,
		fileSuffix string) (*entity.File, error)
	UploadFromUrl(ctx context.Context, userId int, url string, fileType int) (*entity.File, error)
	List(ctx context.Context, uid int, page, size int) (*response.GetFilesResp, error)
	Detail(ctx context.Context, uid int, id int) (*response.FileDetailResp, error)
}

type logic struct {
	model model.Factory
	cache cache.Factory
	srv   service.Factory
}

func NewLogic(model model.Factory, cache cache.Factory, srv service.Factory) Logic {
	return &logic{
		model: model,
		cache: cache,
		srv:   srv,
	}
}

// UploadFile 将用户上传的所有文件在本地进行缓存
func (l *logic) UploadFile(ctx context.Context, userId int, fileContent []byte, fileType int, fileOriginName string,
	fileSuffix string) (*entity.File, error) {
	filename, md5Hash := l.GenerateFilenameByHash(ctx, fileContent, fileSuffix)

	// Using Cloudflare R2
	r2Bucket := "aicbpng" // Hardcode or from config
	contentType := "application/octet-stream"
	if fileSuffix == "png" {
		contentType = "image/png"
	} else if fileSuffix == "jpg" || fileSuffix == "jpeg" {
		contentType = "image/jpeg"
	}

	url, err := l.srv.Cloudflare().UploadToR2(ctx, r2Bucket, filename, fileContent, contentType)
	if err != nil {
		return nil, err
	}
	fileData := &entity.File{
		UserId:       userId,
		Url:          url,
		FileKey:      filename,
		OriginalName: fileOriginName,
		Md5:          md5Hash,
		Type:         fileType,
		Suffix:       fileSuffix,
		Size:         len(fileContent),
		Status:       1,
		CreatedAt:    time.Now(),
	}
	if err := l.model.File().Create(ctx, fileData); err != nil {
		return nil, err
	}
	return fileData, nil
}

func (l *logic) GetFileSuffix(ctx context.Context, filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func (l *logic) GenerateFilenameByHash(ctx context.Context, fileContent []byte, fileSuffix string) (string, string) {
	fileHash := md5.Sum(fileContent)
	md5Hash := hex.EncodeToString(fileHash[:])
	return md5Hash + "." + fileSuffix, md5Hash
}

func (l *logic) List(ctx context.Context, uid int, page, size int) (*response.GetFilesResp, error) {
	offset := (page - 1) * size
	total, files, err := l.model.File().GetList(ctx, uid, offset, size)
	if err != nil {
		return nil, err
	}

	list := make([]*response.FileSummary, 0, len(files))
	for _, f := range files {
		list = append(list, &response.FileSummary{
			ID:        f.Id,
			Url:       f.Url,
			Type:      f.Type,
			CreatedAt: f.CreatedAt,
			Name:      f.OriginalName,
		})
	}

	return &response.GetFilesResp{
		List:  list,
		Total: total,
	}, nil
}

func (l *logic) Detail(ctx context.Context, uid int, id int) (*response.FileDetailResp, error) {
	file, err := l.model.File().GetByIdAndUserId(ctx, id, uid)
	if err != nil {
		return nil, err
	}
	return &response.FileDetailResp{
		File: file,
	}, nil
}

func (l *logic) UploadFromUrl(ctx context.Context, userId int, url string, fileType int) (*entity.File, error) {
	// 1. Download Content
	// Simple HTTP GET. In prod use context aware client with timeout
	resp, err := http.Get(url) // Note: context not used here for simplicity but should be
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 2. Determine file info
	// Prefer extracting filename from URL or header, fall back to hash
	filename := "remote_file"
	// Try parsing url path
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		filename = parts[len(parts)-1]
		// Remove query params
		if idx := strings.Index(filename, "?"); idx != -1 {
			filename = filename[:idx]
		}
	}
	suffix := l.GetFileSuffix(ctx, filename)
	if suffix == "" {
		suffix = "png" // Default?
	}

	// 3. Reuse UploadFile logic
	return l.UploadFile(ctx, userId, content, fileType, filename, suffix)
}
