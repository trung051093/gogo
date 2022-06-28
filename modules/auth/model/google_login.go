package authmodel

type GoogleLogin struct {
	Redirect string `validate:"required" json:"redirect"`
}
