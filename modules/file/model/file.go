package filemodel

import (
	"fmt"
	"strings"
	"user_management/common"
)

const EntityName = "file"

type File struct {
	common.SQLModel
	Url        string `json:"url" gorm:"column:url;"`
	FileName   string `json:"fileName,omitempty" gorm:"column:file_name;"`
	ObjectName string `json:"objectName,omitempty" gorm:"column:object_name;"`
	FilePath   string `json:"filePath,omitempty" gorm:"column:file_path;"`
	MimeType   string `json:"mimeType,omitempty" gorm:"column:mime_type;"`
}

func (File) TableName() string { return "files" }

func (File) TableIndex() string {
	return strings.ToLower(fmt.Sprintf("%s-%s", common.Project, File{}.TableName()))
}
