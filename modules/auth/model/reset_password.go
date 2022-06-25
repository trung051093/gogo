package authmodel

type AuthResetPassword struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
	Token    string `validate:"required" json:"token"`
}
