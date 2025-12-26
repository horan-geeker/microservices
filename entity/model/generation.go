package model

import "time"

const (
	GenerationStatusPending = 0
	GenerationStatusSuccess = 1
	GenerationStatusFailed  = 2
)

type Generation struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID         int       `gorm:"index;column:user_id;comment:关联用户表ID" json:"userId"`
	TraceID        string    `gorm:"column:trace_id;type:varchar(100);comment:外部跟踪ID" json:"traceId"`
	Type           string    `gorm:"column:type;type:varchar(20);comment:内容类型" json:"type"` // text, image, code, audio
	ModelName      string    `gorm:"column:model_name;type:varchar(50);comment:使用的模型" json:"modelName"`
	Prompt         string    `gorm:"column:prompt;type:text;comment:用户原始指令" json:"prompt"`
	NegativePrompt string    `gorm:"column:negative_prompt;type:text;comment:负向提示词" json:"negativePrompt"`
	ContentText    string    `gorm:"column:content_text;type:longtext;comment:生成的内容" json:"contentText"`
	MediaURL       string    `gorm:"column:media_url;type:varchar(500);comment:媒体链接" json:"mediaUrl"`
	Metadata       string    `gorm:"column:metadata;type:json;comment:模型元数据" json:"metadata"` // JSON string
	Status         int       `gorm:"column:status;type:tinyint;default:0;comment:状态" json:"status"`
	IsPublic       bool      `gorm:"column:is_public;default:false;comment:是否公开" json:"isPublic"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Generation) TableName() string {
	return "generations"
}
