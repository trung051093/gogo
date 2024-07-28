package dto

type AuthLogoutDto struct {
	Session string `json:"token" validate:"required" `
}

func (dt *AuthLogoutDto) Validate() error {
	return nil
}
