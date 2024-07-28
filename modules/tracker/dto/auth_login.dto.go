package dto

type AuthLoginDto struct {
	Username string `json:"username" validate:"required" `
	Password string `validate:"required" json:"password"`
}

func (dt *AuthLoginDto) Validate() error {
	return nil
}
