package api

import "golang.org/x/oauth2"

// Author interface is wrapper of api
type Author interface {
	ReadFile(filename string) ([]byte, error)
	ConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error)
	GetToken(config *oauth2.Config, fileName string) *oauth2.Token
}
