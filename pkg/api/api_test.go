package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestService(t *testing.T) {
	setup := func() (*Service, *afero.Afero) {
		fsBackup := AppFs
		AppFs = afero.NewMemMapFs()

		afs := &afero.Afero{Fs: AppFs}

		helperBackup := Helper
		t.Cleanup(func() {
			AppFs = fsBackup
			Helper = helperBackup
		})

		service := &Service{}
		return service, afs
	}

	// language=json
	const credential = `{
  "installed": {
    "client_id": "client-1234",
    "project_id": "project-5678",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "this is sparta!",
    "redirect_uris": [
      "http://localhost"
    ]
  }
}`

	t.Run("success to initialize", func(t *testing.T) {
		service, afs := setup()

		_, caller, _, _ := runtime.Caller(1)
		basePath := filepath.Dir(caller)

		f, err := afs.Create(filepath.Join(basePath, credentialFile))
		assert.NoError(t, err)

		_, err = f.Write([]byte(credential))
		assert.NoError(t, err)

		Helper.GetToken = func(config *oauth2.Config, fileName string) (token *oauth2.Token, err error) {
			return &oauth2.Token{}, nil
		}

		err = service.Init()
		assert.NoError(t, err)
	})

	t.Run("fail to read credential file", func(t *testing.T) {
		service, _ := setup()

		err := service.Init()
		assert.Error(t, err)
	})

	t.Run("handle credentials.json reading error", func(t *testing.T) {
		service, afs := setup()

		_, caller, _, _ := runtime.Caller(1)
		basePath := filepath.Dir(caller)

		f, err := afs.Create(filepath.Join(basePath, credentialFile))
		assert.NoError(t, err)

		given := []byte("invalid json")
		_, err = f.Write(given)
		assert.NoError(t, err)

		err = f.Close()
		assert.NoError(t, err)

		var cfg oauth2.Config
		want := json.Unmarshal(given, &cfg)
		assert.Error(t, want)

		got := service.Init()
		assert.Error(t, got)
		assert.EqualError(t, got, want.Error())
	})

	t.Run("GetToken error", func(t *testing.T) {
		service, afs := setup()

		_, caller, _, _ := runtime.Caller(1)
		basePath := filepath.Dir(caller)

		f, err := afs.Create(filepath.Join(basePath, credentialFile))
		assert.NoError(t, err)

		_, err = f.Write([]byte(credential))
		assert.NoError(t, err)

		given := errors.New("get token error")
		Helper.GetToken = func(config *oauth2.Config, fileName string) (*oauth2.Token, error) {
			return nil, given
		}

		err = service.Init()
		assert.Error(t, err)
		assert.EqualError(t, err, given.Error())
	})
}

func TestOAuthCallbackServer(t *testing.T) {
	const given = "want"
	want := given

	for _, tc := range []struct {
		description string
		query       string
		success     bool
		err         error
	}{
		{
			description: "success",
			query:       fmt.Sprintf("?code=%s", want),
			success:     true,
			err:         nil,
		},
		{
			description: "empty param",
			query:       "",
			success:     false,
			err:         errors.New("invalid callback params \n"),
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			ln, err := net.Listen("tcp", ":0")
			assert.NoError(t, err)

			addr := strings.Split(ln.Addr().String(), "]")[1]
			go func() {
				time.Sleep(500 * time.Millisecond)
				reqUrl := fmt.Sprintf("http://localhost%s/auth/callback/%s", addr, tc.query)
				req, err := http.NewRequest("GET", reqUrl, nil)
				assert.NoError(t, err)

				client := &http.Client{}
				resp, err := client.Do(req)
				assert.NoError(t, err)

				defer func() {
					err := resp.Body.Close()
					assert.NoError(t, err)
				}()
			}()

			actual, err := GetTokenWithBrowser(ln)
			if tc.success {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.success, actual == want)
		})
	}
}
