package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// NewService function returns initialized Service object's pointer
func NewService() (*Service, error) {
	api := &Service{
		AuthImpl: &AuthImpl{
			ConfigFromJSON: google.ConfigFromJSON,
			ReadFile:       ioutil.ReadFile,
			GetToken:       GetToken,
		},
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

	*AuthImpl
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

	token, err := api.GetToken(config, tokenFile, Func{
		TokenFromFile: TokenFromFile,
		TokenFromWeb:  TokenFromWeb,
		GetAuthCode: func() (s string, e error) {
			_, e = fmt.Scan(&s)
			return
		},
		SaveToken: SaveToken,
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

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		//noinspection SpellCheckingInspection
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

// TokenFromWeb request a token from the web, then returns the retrieved token.
func TokenFromWeb(config *oauth2.Config, getAuthCode func() (string, error)) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n\n>>", authURL)
	if err := openBrowser(authURL); err != nil {
		log.Fatal(err)
	}

	authCode, err := getAuthCode()
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
