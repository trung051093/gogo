package authmodel

type AuthLogin struct {
	Email    string `validate:"required,email" json:"email" gorm:"-"`
	Password string `validate:"required" json:"password" gorm:"-"`
}
