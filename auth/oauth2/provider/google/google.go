package google

import (
	"encoding/json"

	"github.com/core-stack/authz/auth/oauth2"
)

func Provider(clientId, secret, redirectUrl string, scopes ...string) oauth2.Provider {
	if scopes == nil || scopes[0] == "" {
		scopes = []string{"email", "profile"}
	}
	return oauth2.Provider{
		Name:         "google",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
		Scopes:       scopes,
		ClientID:     clientId,
		ClientSecret: secret,
		UserInfoURL:  "https://www.googleapis.com/oauth2/v3/userinfo",
		RedirectURI:  redirectUrl,
		ExtractUserInfo: func(data []byte) (oauth2.UserInfo, error) {
			var userInfo oauth2.UserInfo
			if err := json.Unmarshal(data, &userInfo); err != nil {
				return oauth2.UserInfo{}, err
			}
			return userInfo, nil
		},
	}
}
