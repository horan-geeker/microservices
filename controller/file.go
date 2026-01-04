package controller

import (
	"io"
	"microservices/cache"
	"microservices/entity/consts"
	"microservices/entity/ecode"
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileController interface {
	Upload(c *gin.Context, param *request.UploadFileParam) (*response.UploadFile, error)
	List(c *gin.Context, req *request.GetFilesReq) (*response.GetFilesResp, error)
	Detail(c *gin.Context, id int) (*response.FileDetailResp, error)
}

type fileController struct {
	logic   logic.Factory
	service service.Factory
}

// NewFileController .
func NewFileController(model model.Factory, cache cache.Factory, service service.Factory) FileController {
	return &fileController{
		logic:   logic.NewLogic(model, cache, service),
		service: service,
	}
}

// Upload godoc
// @Summary 上传文件
// @Description 上传图片文件，限制为图片格式，每秒只允许执行一次
// @Tags files
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "图片文件"
// @Success 200 {object} entity.Response[response.UploadFile]
// @Failure 400 {object} entity.Response[any]
// @Failure 401 {object} entity.Response[any] "用户登录态校验失败"
// @Failure 429 {object} entity.Response[any] "请求过于频繁"
// @Router /files/upload [post]
func (f *fileController) Upload(c *gin.Context, param *request.UploadFileParam) (*response.UploadFile, error) {
	ctx := c.Request.Context()
	// 获取当前登录用户
	auth, err := f.logic.Auth().GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	user, err := f.logic.User().GetByUid(ctx, auth.Uid)
	if err != nil {
		return nil, err
	}

	// 获取上传的文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, ecode.ErrParamInvalid.WithMessage("请选择要上传的文件")
	}

	// 验证文件扩展名是否为图片格式
	filename := fileHeader.Filename
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return nil, ecode.ErrParamInvalid.WithMessage("只支持图片格式文件 (jpg, jpeg, png, gif, bmp, webp)")
	}

	// 读取文件内容
	file, err := fileHeader.Open()
	if err != nil {
		return nil, ecode.ErrParamInvalid.WithMessage("无法读取上传的文件")
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, ecode.ErrParamInvalid.WithMessage("无法读取文件内容")
	}

	// 验证文件大小 (限制为10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if len(fileContent) > maxSize {
		return nil, ecode.ErrParamInvalid.WithMessage("文件大小不能超过10MB")
	}

	// 生成新的文件名
	suffix := strings.TrimPrefix(ext, ".")

	// 上传文件到云存储
	fileData, err := f.logic.File().UploadFile(ctx, user.ID, fileContent, consts.FileTypeImage, filename, suffix)
	if err != nil {
		return nil, err
	}

	return &response.UploadFile{
		File: *fileData,
	}, nil
}

// List godoc
// @Summary 获取文件列表
// @Description 分页获取用户上传的文件列表
// @Tags files
// @Accept  json
// @Produce  json
// @Param page query int true "页码"
// @Param size query int true "每页数量"
// @Success 200 {object} entity.Response[response.GetFilesResp]
// @Router /files [get]
func (f *fileController) List(c *gin.Context, req *request.GetFilesReq) (*response.GetFilesResp, error) {
	auth, err := f.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return f.logic.File().List(c.Request.Context(), auth.Uid, req.Page, req.Size)
}

// Detail godoc
// @Summary 获取文件详情
// @Description 获取指定ID的文件详情
// @Tags files
// @Accept  json
// @Produce  json
// @Param id path int true "文件ID"
// @Success 200 {object} entity.Response[response.FileDetailResp]
// @Router /files/{id} [get]
func (f *fileController) Detail(c *gin.Context, id int) (*response.FileDetailResp, error) {
	auth, err := f.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return f.logic.File().Detail(c.Request.Context(), auth.Uid, id)
}
