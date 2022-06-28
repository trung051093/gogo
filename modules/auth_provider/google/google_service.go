package googleauthprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"gogo/common"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectUri  string
}

type googleAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleAuthProvider(config *GoogleAuthConfig) *googleAuthProvider {
	cfg := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectUri,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &googleAuthProvider{config: cfg}
}

func (g *googleAuthProvider) GetAuthUri(ctx context.Context, state string, opts ...oauth2.AuthCodeOption) string {
	url := g.config.AuthCodeURL(state, opts...)
	return url
}

// Use code to get token and get user info from Google.
func (g *googleAuthProvider) GetUser(ctx context.Context, code string) (*GoogleUser, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, common.NewErrorResponse(err, "code exchange wrong", "GoogleAuth", "GoogleAuth")
	}

	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		return nil, common.NewErrorResponse(err, "failed to get user info", "GoogleAuth", "GoogleAuth")
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, common.NewErrorResponse(err, "failed read response", "GoogleAuth", "GoogleAuth")
	}

	var user GoogleUser
	if err := json.Unmarshal(contents, &user); err != nil {
		return nil, common.NewErrorResponse(err, "failed to parse user info", "GoogleAuth", "GoogleAuth")
	}

	return &user, nil
}
