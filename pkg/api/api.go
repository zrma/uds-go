package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/zrma/uds-go/pkg/api/browser"
	"github.com/zrma/uds-go/pkg/uds"
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
	_, caller, _, _ := runtime.Caller(2)
	basePath := filepath.Dir(caller)

	b, err := api.ReadFile(filepath.Join(basePath, credentialFile))
	if err != nil {
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := api.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return err
	}

	token, err := api.GetToken(config, filepath.Join(basePath, tokenFile), Func{
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

// GetBaseFolder locate the base UDS folder
func (api *Service) GetBaseFolder() (*drive.File, error) {
	r, err := api.Files.List().
		Q("properties has {key='udsRoot' and value='true'} and trashed=false").
		PageSize(1).
		Fields("nextPageToken, files(id, name, properties)").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	fileLength := len(r.Files)
	if fileLength == 0 {
		fmt.Println("No files found.")
		return api.createRootFolder()
	} else if fileLength == 1 {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
		return r.Files[0], nil
	}
	return nil, fmt.Errorf("multiple UDS Roots found")
}

func (api *Service) createRootFolder() (*drive.File, error) {
	return api.Files.Create(&drive.File{
		Name:       "UDS Root",
		MimeType:   "application/vnd.google-apps.folder",
		Properties: map[string]string{"udsRoot": "true"},
		Parents:    []string{},
	}).Fields("id").Do()
}

func (api *Service) CreateMediaFolder(media *uds.File) (*drive.File, error) {
	return api.Files.Create(&drive.File{
		Name:     media.Name,
		MimeType: "application/vnd.google-apps.folder",
		Properties: map[string]string{
			"udsRoot":      "true",
			"size":         media.Size,
			"size_numeric": media.SizeNumeric,
			"encoded_size": media.EncodedSize,
			"md5":          media.MD5,
		},
		Parents: media.Parents,
	}).Fields("id").Do()
}

func (api *Service) ListFiles(query string) ([]*uds.File, error) {
	q := "properties has {key='uds' and value='true'} and trashed=false"
	if query != "" {
		q += fmt.Sprintf(" and name contains '%s'", query)
	}

	r, err := api.Files.List().
		Q(q).
		PageSize(1000).
		Fields("nextPageToken, files(id, name, properties, mimeType)").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	var files []*uds.File
	for _, f := range r.Files {
		props := f.Properties
		file := &uds.File{
			Name:        props["name"],
			Mime:        props["mimeType"],
			Size:        props["size"],
			EncodedSize: props["encoded_size"],
			SizeNumeric: props["size_numeric"],
			ID:          props["id"],
			MD5:         props["md5Checksum"],
			Shared:      props["shared"] == "true",
		}
		file.Init()
		files = append(files, file)
	}
	return files, nil
}

// GetTokenWithBrowser function receive token with localhost callback server
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

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println("callback listen server shutdown", err)
	}
	return authCodes[0], nil
}

func getAuthCodeOffline(config *oauth2.Config, f Func) (string, error) {
	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n\n>>", authURL)
	return f.GetAuthCode()
}

func getAuthCode(config *oauth2.Config, f Func) (string, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	defer ln.Close()

	addr := strings.Split(ln.Addr().String(), "]")[1]
	config.RedirectURL = fmt.Sprintf("http://localhost%s/auth/callback/", addr)
	if f.OpenBrowser == nil {
		f.OpenBrowser = func(s string) error {
			return errors.New("impossible to open a browser")
		}
	}
	var authCode string
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if err := f.OpenBrowser(authURL); err != nil {
		return getAuthCodeOffline(config, f)
	} else if authCode, err = GetTokenWithBrowser(ln); err != nil {
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
