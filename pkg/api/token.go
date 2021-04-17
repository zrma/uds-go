package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/spf13/afero"
	context2 "golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/zrma/uds-go/pkg/api/browser"
)

type helper struct {
	ConfigFromJSON   func(jsonKey []byte, scope ...string) (*oauth2.Config, error)
	GetTokenFromFile func(filePath string) (*oauth2.Token, error)
	GetTokenFromWeb  func(config *oauth2.Config) (*oauth2.Token, error)
	ScanAuthCode     func() (string, error)
	ExchangeToken    func(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error)
	OpenBrowser      func(url string) error
	GetToken         func(config *oauth2.Config, fileName string) (*oauth2.Token, error)
}

var Helper helper

func init() {
	Helper = helper{
		ConfigFromJSON:   google.ConfigFromJSON,
		GetTokenFromFile: getTokenFromFile,
		GetTokenFromWeb:  getTokenFromWeb,
		ScanAuthCode: func() (s string, e error) {
			_, e = fmt.Scan(&s)
			return
		},
		OpenBrowser: browser.Open,
		GetToken:    getToken,
	}
}

func getToken(config *oauth2.Config, fileName string) (*oauth2.Token, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	token, err := Helper.GetTokenFromFile(fileName)
	if err != nil {
		token, err = Helper.GetTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err = saveToken(fileName, token); err != nil {
			return nil, err
		}
	}
	if !token.Valid() {
		if Helper.ExchangeToken != nil {
			return Helper.ExchangeToken(config, token)
		}
		token, err = config.TokenSource(context.Background(), token).Token()
		if err != nil {
			return nil, err
		}
		if err := saveToken(fileName, token); err != nil {
			return nil, err
		}
		return token, nil
	}
	return token, nil
}

const (
	credentialFile = "credentials.json"
	tokenFile      = "token.json"
)

func GetTokenWithBrowser(ln net.Listener) (string, error) {
	tokenCh := make(chan string)
	defer close(tokenCh)

	var do sync.Once

	handler := http.NewServeMux()
	handler.HandleFunc("/auth/callback/", func(w http.ResponseWriter, r *http.Request) {
		const response = `
<body>
Finished
</body>
`
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
		do.Do(func() {
			tokenCh <- r.URL.RawQuery
		})
	})

	if ln == nil {
		var err error
		ln, err = net.Listen("tcp", ":0")
		if err != nil {
			return "", err
		}
	}

	addr := strings.Split(ln.Addr().String(), "]")[1]
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			log.Println("callback listen server closed", err)
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

	if err := srv.Shutdown(context2.Background()); err != nil {
		log.Println("callback listen server shutdown", err)
	}
	return authCodes[0], nil
}

func getAuthCodeOffline(config *oauth2.Config) (string, error) {
	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n\n>>", authURL)
	return Helper.ScanAuthCode()
}

func getAuthCode(config *oauth2.Config) (authCode string, err error) {
	var ln net.Listener
	ln, err = net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	defer func() {
		err = ln.Close()
	}()

	addr := strings.Split(ln.Addr().String(), "]")[1]
	config.RedirectURL = fmt.Sprintf("http://localhost%s/auth/callback/", addr)
	if Helper.OpenBrowser == nil {
		return getAuthCodeOffline(config)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if err := Helper.OpenBrowser(authURL); err != nil {
		return getAuthCodeOffline(config)
	}

	if authCode, err = GetTokenWithBrowser(ln); err != nil {
		return getAuthCodeOffline(config)
	}

	return authCode, nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authCode, err := getAuthCode(config)
	if err != nil {
		return nil, err
	}

	fmt.Println(authCode)
	token, err := config.Exchange(context2.Background(), authCode)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getTokenFromFile(file string) (token *oauth2.Token, err error) {
	var f afero.File
	f, err = AppFs.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()
	token = &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func saveToken(path string, token *oauth2.Token) (err error) {
	fmt.Printf("Saving credential file to: %s\n", path)
	var f afero.File
	f, err = AppFs.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	return json.NewEncoder(f).Encode(token)
}
