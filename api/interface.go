package api

import "golang.org/x/oauth2"

type Reader interface {
	ReadFile(filename string) ([]byte, error)
	ConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error)
}

type Author interface {
	GetToken(config *oauth2.Config) *oauth2.Token
}
