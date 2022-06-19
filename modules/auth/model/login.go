package authmodel

import authprovider "user_management/modules/auth_providers"

type AuthLogin struct {
	Email    string `validate:"required,email" json:"email" gorm:"-"`
	Password string `validate:"required" json:"password" gorm:"-"`
}

type AuthResponse struct {
	authprovider.TokenProvider
}
