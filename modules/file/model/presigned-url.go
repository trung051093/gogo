package filemodel

type PresignedPostObject struct {
	Url    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}
