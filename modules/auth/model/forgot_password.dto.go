package authmodel

type AuthForgotPasswordDto struct {
	Email             string `validate:"required,email" json:"email"`
	ForgotPasswordUri string `validate:"required" json:"forgotPasswordUri"`
}
