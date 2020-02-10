package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/zrma/uds-go/pkg/api/browser"
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
		OpenBrowser:   browser.Open,
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

func getTokenWithBrowser() (string, error) {
	tokenCh := make(chan string)
	handler := http.NewServeMux()
	handler.HandleFunc("/auth/callback/", func(w http.ResponseWriter, r *http.Request) {
		const response = `
<body>
Finished
</body>
`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		tokenCh <- r.URL.RawQuery
	})
	srv := &http.Server{
		Addr:    ":1333",
		Handler: handler,
	}
	go func() {
		if err := srv.ListenAndServe(); err != srv.ListenAndServe() {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	rawQuery := <-tokenCh
	m, err := url.ParseQuery(rawQuery)
	if err != nil {
		return "", err
	}

	authCodes, ok := m["code"]
	if !ok || len(authCodes) == 0 {
		return "", errors.New(fmt.Sprintln("invalid callback params", rawQuery))
	}

	go func() {
		defer close(tokenCh)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("srv.Shutdown(): %v", err)
		}
	}()
	return authCodes[0], nil
}

func getAuthCodeOffline(config *oauth2.Config, f Func) (string, error) {
	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n\n>>", authURL)
	return f.GetAuthCode()
}

func getAuthCode(config *oauth2.Config, f Func) (string, error) {
	config.RedirectURL = "http://localhost:1333/auth/callback/"
	if f.OpenBrowser == nil {
		f.OpenBrowser = func(s string) error {
			return errors.New("impossible to open a browser")
		}
	}
	var authCode string
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if err := f.OpenBrowser(authURL); err != nil {
		return getAuthCodeOffline(config, f)
	} else if authCode, err = getTokenWithBrowser(); err != nil {
		return getAuthCodeOffline(config, f)
	}
	return authCode, nil
}

// TokenFromWeb request a token from the web, then returns the retrieved token.
func TokenFromWeb(config *oauth2.Config, f Func) (*oauth2.Token, error) {
	authCode, err := getAuthCode(config, f)
	if err != nil {
		return nil, err
	}

	fmt.Println(authCode)
	token, err := config.Exchange(context.Background(), authCode)
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
