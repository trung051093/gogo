package authmodel

type AuthForgotPassword struct {
	Email             string `validate:"required,email" json:"email" gorm:"-"`
	ForgotPasswordUri string `validate:"required" json:"forgotPasswordUri" gorm:"-"`
}
