package authmodel

type AuthRegister struct {
	Email     string `validate:"required,email" json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `validate:"required" json:"password"`
}
