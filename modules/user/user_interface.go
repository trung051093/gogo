package user

import (
	"context"
	"user_management/common"
	usermodel "user_management/modules/user/model"
)

//Reader interface
type Reader interface {
	Get(ctx context.Context, id uint) (*usermodel.User, error)
	Search(ctx context.Context, cond map[string]interface{}) ([]usermodel.User, error)
}

//Writer user writer
type Writer interface {
	Create(ctx context.Context, user *usermodel.User) (int, error)
	Update(ctx context.Context, id uint, userUpdate *usermodel.UserUpdate) error
	Delete(ctx context.Context, user *usermodel.User) error
}

//Repository interface
type UserRepository interface {
	Reader
	Writer
}

//Service interface
type UserService interface {
	GetUser(ctx context.Context, id uint) (*usermodel.User, error)
	SearchUser(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error)
	SearchUsers(ctx context.Context, cond map[string]interface{}, filter *usermodel.UserFilter, paging *common.Pagination) ([]usermodel.User, error)
	CreateUser(ctx context.Context, user *usermodel.UserCreate) (int, error)
	UpdateUser(ctx context.Context, id uint, userUpdate *usermodel.UserUpdate) error
	DeleteUser(ctx context.Context, id uint) error
}
