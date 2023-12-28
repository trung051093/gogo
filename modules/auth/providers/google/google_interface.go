package googleauthprovider

import (
	"context"

	"golang.org/x/oauth2"
)

type GoogleAuthProvider interface {
	GetAuthUri(ctx context.Context, state string, opts ...oauth2.AuthCodeOption) string
	GetUser(ctx context.Context, code string) (*GoogleUser, error)
	GetUserTrace(ctx context.Context, code string) (*GoogleUser, error)
}
