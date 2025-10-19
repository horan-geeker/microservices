package model

import "time"

type File struct {
	Id           int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	UserId       int       `gorm:"column:user_id;type:int(11);NOT NULL" json:"user_id"`
	Url          string    `gorm:"column:url;type:varchar(500);NOT NULL" json:"url"`
	FileKey      string    `gorm:"column:file_key;type:varchar(255);NOT NULL" json:"file_key"`
	OriginalName string    `gorm:"column:original_name;type:varchar(255);NOT NULL" json:"original_name"`
	Md5          string    `gorm:"column:md5;type:varchar(255);NOT NULL" json:"md5"`
	Type         int       `gorm:"column:type;type:int(11);NOT NULL" json:"type"`
	Suffix       string    `gorm:"column:suffix;type:varchar(255);NOT NULL" json:"suffix"`
	Size         int       `gorm:"column:size;type:int(11);NOT NULL" json:"size"`
	Status       int       `gorm:"column:status;type:int(11);NOT NULL" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;NOT NULL" json:"created_at"`
	PublicUrl    string    `gorm:"-" json:"public_url"`
}
