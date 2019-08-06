package api

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
)

// Service struct is google api service wrapper.
type Service struct {
	*drive.Service
	ctx context.Context

	Author
}

const (
	credentialFile = "credentials.json"
	tokenFile      = "token.json"
)

// NewService function returns initialized Service object's pointer
func NewService() (*Service, error) {
	api := &Service{
		Author: &AuthorImpl{},
	}
	if err := api.Init(); err != nil {
		return nil, err
	}
	return api, nil
}

// Init works internally but public(export) for using in apt_test package
func (api *Service) Init() error {
	b, err := api.ReadFile(credentialFile)
	if err != nil {
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := api.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return err
	}

	token, err := api.GetToken(config, tokenFile, func() (s string, e error) {
		_, e = fmt.Scan(&s)
		return
	})
	if err != nil {
		return err
	}

	api.ctx = context.Background()
	driveService, err := drive.NewService(
		api.ctx,
		option.WithTokenSource(config.TokenSource(api.ctx, token)),
	)
	if err != nil {
		return err
	}

	api.Service = driveService
	return nil
}

// AuthorImpl is impl some api for mocking test
type AuthorImpl struct {
}

// ConfigFromJSON wrapping google.ConfigFromJSON api
func (AuthorImpl) ConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error) {
	return google.ConfigFromJSON(jsonKey, scope...)
}

// ReadFile wrapping ioutil.ReadFile
func (AuthorImpl) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// GetToken return oauth token
func (a AuthorImpl) GetToken(config *oauth2.Config, fileName string, f func() (string, error)) (*oauth2.Token, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	token, err := TokenFromFile(fileName)
	if err != nil {
		token, err = GetTokenFromWeb(config, f)
		if err != nil {
			return nil, err
		}
		if err := SaveToken(fileName, token); err != nil {
			return nil, err
		}
	}
	return token, nil
}

// GetTokenFromWeb request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config, f func() (string, error)) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n\n>>", authURL)

	authCode, err := f()
	if err != nil {
		return nil, err
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TokenFromFile retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// SaveToken saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
