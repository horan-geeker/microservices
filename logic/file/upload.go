package file

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"microservices/cache"
	entity "microservices/entity/model"
	"microservices/model"
	"microservices/service"
	"strings"
	"time"
)

type Logic interface {
	SendSmsCode(ctx context.Context, phone string, code string) error
}

type logic struct {
	model model.Factory
	cache cache.Factory
	srv   service.Factory
}

// UploadFile 将用户上传的所有文件在本地进行缓存
func (l *logic) UploadFile(ctx context.Context, userId int, file io.Reader, fileType int, fileOriginName string,
	fileSuffix string) (*entity.File, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	filename, md5Hash := l.GenerateFilenameByHash(ctx, fileContent, fileSuffix)
	url, fileKey, err := l.srv.Tencent().UploadToCOS(ctx, userId, fileContent, filename)
	if err != nil {
		return nil, err
	}
	fileData := &entity.File{
		UserId:       userId,
		Url:          url,
		FileKey:      fileKey,
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
