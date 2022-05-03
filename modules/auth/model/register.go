package authmodel

type AuthRegister struct {
	Email    string `validate:"required,email" json:"email" gorm:"-"`
	Password string `validate:"required" json:"password" gorm:"-"`
}
