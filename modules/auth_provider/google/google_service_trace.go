package googleauthprovider

import (
	"context"
	decorator "gogo/decorators"
)

func (g *googleAuthProvider) GetUserTrace(ctx context.Context, code string) (*GoogleUser, error) {
	data, err := decorator.TraceService[*GoogleUser](ctx, "googleAuthProvider.GetUser")(g, "GetUser")(ctx, code)
	return data, err
}
