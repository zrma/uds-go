package api

import (
	"golang.org/x/oauth2"
)

// AuthImpl is implementing some api for mocking test
type AuthImpl struct {
	ConfigFromJSON func(jsonKey []byte, scope ...string) (*oauth2.Config, error)
	ReadFile       func(filename string) ([]byte, error)
	GetToken       func(config *oauth2.Config, fileName string, f Func) (*oauth2.Token, error)
}

// Func is internal factory functions to call with dependency
type Func struct {
	TokenFromFile func(file string) (*oauth2.Token, error)
	TokenFromWeb  func(config *oauth2.Config, getAuthCode func() (string, error)) (*oauth2.Token, error)
	GetAuthCode   func() (string, error)
	SaveToken     func(path string, token *oauth2.Token) error
}

// GetToken return oauth token
func GetToken(config *oauth2.Config, fileName string, f Func) (*oauth2.Token, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	token, err := f.TokenFromFile(fileName)
	if err != nil || !token.Valid() {
		token, err = f.TokenFromWeb(config, f.GetAuthCode)
		if err != nil {
			return nil, err
		}
		if err := f.SaveToken(fileName, token); err != nil {
			return nil, err
		}
	}
	return token, nil
}
