package request

type GetFilesReq struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=100"`
}

type FileUpload struct {
	// Might be used later if we need strict struct for file upload params
}

type UploadFileParam struct {
	// Just placeholder if binding needed
}
