package user

import usermodel "user_management/modules/user/model"

//Reader interface
type Reader interface {
	Get(id uint) (*usermodel.User, error)
	Search(cond map[string]interface{}) ([]usermodel.User, error)
}

//Writer user writer
type Writer interface {
	Create(user *usermodel.User) error
	Update(cond map[string]interface{}, userUpdate *usermodel.UserUpdate) error
	Delete(cond map[string]interface{}) (*usermodel.User, error)
}

//Repository interface
type UserRepository interface {
	Reader
	Writer
}

//Service interface
type UserService interface {
	GetUser(id uint) (*usermodel.User, error)
	SearchUsers(cond map[string]interface{}) ([]usermodel.User, error)
	CreateUser(user *usermodel.User) error
	UpdateUser(cond map[string]interface{}, userUpdate *usermodel.UserUpdate) error
	DeleteUser(cond map[string]interface{}) (*usermodel.User, error)
}
