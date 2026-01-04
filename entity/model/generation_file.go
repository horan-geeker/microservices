package model

import "time"

const (
	GenerationFileTypeInput  = 1
	GenerationFileTypeOutput = 2
)

type GenerationFile struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	GenerationID uint64    `gorm:"index;column:generation_id;comment:生成记录ID" json:"generationId"`
	FileID       int       `gorm:"index;column:file_id;comment:文件ID" json:"fileId"`
	Type         int       `gorm:"column:type;comment:文件类型(1:原图,2:生成图)" json:"type"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (GenerationFile) TableName() string {
	return "generation_files"
}
