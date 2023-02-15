package user

import (
	"context"
	"gogo/common"
	elasticsearchmodel "gogo/components/elasticsearch/model"
	decorator "gogo/decorators"
	usermodel "gogo/modules/user/model"
)

func (s *userService) SearchUsersTrace(ctx context.Context, cond map[string]interface{}, f *usermodel.UserFilter, p *common.PagePagination) ([]usermodel.User, error) {
	data, err := decorator.TraceService[[]usermodel.User](ctx, "userService.SearchUsers")(s, "SearchUsers")(ctx, cond, f, p)
	return data, err
}

func (s *userService) SearchUserTrace(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error) {
	data, err := decorator.TraceService[*usermodel.User](ctx, "userService.SearchUser")(s, "SearchUser")(ctx, cond)
	return data, err
}

func (s *userService) GetUserTrace(ctx context.Context, id int) (*usermodel.User, error) {
	data, err := decorator.TraceService[*usermodel.User](ctx, "userService.GetUser")(s, "GetUser")(ctx, id)
	return data, err
}

func (s *userService) CreateUserTrace(ctx context.Context, newUser *usermodel.UserCreate) (int, error) {
	data, err := decorator.TraceService[int](ctx, "userService.CreateUser")(s, "CreateUser")(ctx, newUser)
	return data, err
}

func (s *userService) UpdateUserTrace(ctx context.Context, id int, userUpdate *usermodel.UserUpdate) (int, error) {
	data, err := decorator.TraceService[int](ctx, "userService.UpdateUser")(s, "UpdateUser")(ctx, id, userUpdate)
	return data, err
}

func (s *userService) DeleteUserTrace(ctx context.Context, id int) (int, error) {
	data, err := decorator.TraceService[int](ctx, "userService.DeleteUser")(s, "DeleteUser")(ctx, id)
	return data, err
}

func (s *userService) EsSearchTrace(ctx context.Context, query string, lastIndex string, f *usermodel.UserFilter, p *common.PagePagination) (*elasticsearchmodel.SearchResults, error) {
	data, err := decorator.TraceService[*elasticsearchmodel.SearchResults](ctx, "userService.EsSearch")(s, "EsSearch")(ctx, query, lastIndex, f, p)
	return data, err
}
