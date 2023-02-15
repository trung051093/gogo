package user

import (
	"context"
	"gogo/common"
	elasticsearchmodel "gogo/components/elasticsearch/model"
	usermodel "gogo/modules/user/model"
)

// Reader interface
type Reader interface {
	Get(ctx context.Context, id int) (*usermodel.User, error)
	Search(ctx context.Context, cond map[string]interface{}, filter *usermodel.UserFilter, paging *common.PagePagination) ([]usermodel.User, error)
	SearchOne(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error)
}

// Writer user writer
type Writer interface {
	Create(ctx context.Context, user *usermodel.UserCreate) (int, error)
	Update(ctx context.Context, id int, userUpdate *usermodel.UserUpdate) (int, error)
	DeActive(ctx context.Context, user *usermodel.User) (int, error)
	Delete(ctx context.Context, user *usermodel.User) (int, error)
}

// Repository interface
type UserRepository interface {
	Reader
	Writer
}

type UserServiceTrace interface {
	SearchUsersTrace(ctx context.Context, cond map[string]interface{}, f *usermodel.UserFilter, p *common.PagePagination) ([]usermodel.User, error)
	SearchUserTrace(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error)
	GetUserTrace(ctx context.Context, id int) (*usermodel.User, error)
	CreateUserTrace(ctx context.Context, newUser *usermodel.UserCreate) (int, error)
	UpdateUserTrace(ctx context.Context, id int, userUpdate *usermodel.UserUpdate) (int, error)
	DeleteUserTrace(ctx context.Context, id int) (int, error)
	EsSearchTrace(ctx context.Context, query string, lastIndex string, f *usermodel.UserFilter, p *common.PagePagination) (*elasticsearchmodel.SearchResults, error)
}

//Service interface
type UserService interface {
	// trace api
	UserServiceTrace

	// normal
	GetUser(ctx context.Context, id int) (*usermodel.User, error)
	SearchUser(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error)
	SearchUsers(ctx context.Context, cond map[string]interface{}, filter *usermodel.UserFilter, paging *common.PagePagination) ([]usermodel.User, error)
	CreateUser(ctx context.Context, user *usermodel.UserCreate) (int, error)
	UpdateUser(ctx context.Context, id int, userUpdate *usermodel.UserUpdate) (int, error)
	DeleteUser(ctx context.Context, id int) (int, error)
	EsSearch(ctx context.Context, query string, lastIndex string, f *usermodel.UserFilter, p *common.PagePagination) (*elasticsearchmodel.SearchResults, error)
}
