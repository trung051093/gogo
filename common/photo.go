package common

type Photo struct {
	SQLModel
	Url        string `json:"url,omitempty" gorm:"column:url;"`
	FileName   string `json:"fileName,omitempty" gorm:"column:file_name;"`
	ObjectName string `json:"objectName,omitempty" gorm:"column:object_name;"`
	FilePath   string `json:"filePath,omitempty" gorm:"column:file_path;"`
	MimeType   string `json:"mimeType,omitempty" gorm:"column:mime_type;"`
	Width      string `json:"width,omitempty" gorm:"column:width;"`
	Height     string `json:"height,omitempty" gorm:"column:height;"`
}
