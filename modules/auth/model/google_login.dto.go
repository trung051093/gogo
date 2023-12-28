package authmodel

type GoogleLoginDto struct {
	Redirect string `validate:"required" json:"redirect"`
}
