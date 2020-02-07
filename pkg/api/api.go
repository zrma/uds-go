package api

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// NewService function returns initialized Service object's pointer
func NewService() (*Service, error) {
	api := &Service{
		Auth: &AuthImpl{},
	}
	if err := api.Init(); err != nil {
		return nil, err
	}
	return api, nil
}

// Service struct is google api service wrapper.
type Service struct {
	*drive.Service
	ctx context.Context

	Auth
}

const (
	credentialFile = "credentials.json"
	tokenFile      = "token.json"
)

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
