package authmodel

import authmodelprovider "gogo/modules/auth_provider/model"

type AuthLogin struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

type AuthResponse struct {
	authmodelprovider.TokenProvider
}
