package api

import "golang.org/x/oauth2"

// Auth interface is wrapper of api
type Auth interface {
	ReadFile(filename string) ([]byte, error)
	ConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error)
	GetToken(config *oauth2.Config, fileName string, f func() (string, error)) (*oauth2.Token, error)
}
