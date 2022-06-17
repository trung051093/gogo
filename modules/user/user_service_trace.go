package user

import (
	"context"
	"user_management/common"
	elasticsearchmodel "user_management/components/elasticsearch/model"
	decorator "user_management/decorators"
	usermodel "user_management/modules/user/model"
)

func (s *userService) SearchUsersTrace(ctx context.Context, cond map[string]interface{}, f *usermodel.UserFilter, p *common.Pagination) ([]usermodel.User, error) {
	data, err := decorator.TraceService(ctx, "userService.SearchUsers")(s, "SearchUsers")(ctx, cond, f, p)
	return data.([]usermodel.User), err
}

func (s *userService) SearchUserTrace(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error) {
	data, err := decorator.TraceService(ctx, "userService.SearchUser")(s, "SearchUser")(ctx, cond)
	return data.(*usermodel.User), err
}

func (s *userService) GetUserTrace(ctx context.Context, id int) (*usermodel.User, error) {
	data, err := decorator.TraceService(ctx, "userService.GetUser")(s, "GetUser")(ctx, id)
	return data.(*usermodel.User), err
}

func (s *userService) CreateUserTrace(ctx context.Context, newUser *usermodel.UserCreate) (int, error) {
	data, err := decorator.TraceService(ctx, "userService.CreateUser")(s, "CreateUser")(ctx, newUser)
	return data.(int), err
}

func (s *userService) UpdateUserTrace(ctx context.Context, id int, userUpdate *usermodel.UserUpdate) (int, error) {
	data, err := decorator.TraceService(ctx, "userService.UpdateUser")(s, "UpdateUser")(ctx, id, userUpdate)
	return data.(int), err
}

func (s *userService) DeleteUserTrace(ctx context.Context, id int) (int, error) {
	data, err := decorator.TraceService(ctx, "userService.DeleteUser")(s, "DeleteUser")(ctx, id)
	return data.(int), err
}

func (s *userService) EsSearchTrace(ctx context.Context, query string, lastIndex string, f *usermodel.UserFilter, p *common.Pagination) (*elasticsearchmodel.SearchResults, error) {
	data, err := decorator.TraceService(ctx, "userService.EsSearch")(s, "EsSearch")(ctx, query, lastIndex, f, p)
	return data.(*elasticsearchmodel.SearchResults), err
}
