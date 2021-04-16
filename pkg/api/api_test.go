package api_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"github.com/zrma/uds-go/pkg/api"
)

func TestService(t *testing.T) {
	setup := func() (*api.Service, *api.AuthImpl) {
		authImpl := &api.AuthImpl{
			ConfigFromJSON: func(jsonKey []byte, scope ...string) (config *oauth2.Config, err error) {
				return &oauth2.Config{}, nil
			},
			ReadFile: func(filename string) (bytes []byte, err error) {
				return []byte{}, nil
			},
			GetToken: func(config *oauth2.Config, fileName string, f api.Func) (token *oauth2.Token, err error) {
				return &oauth2.Token{}, nil
			},
		}
		service := &api.Service{
			AuthImpl: authImpl,
		}
		return service, authImpl
	}

	t.Run("success to initialize", func(t *testing.T) {
		service, _ := setup()
		err := service.Init()
		assert.NoError(t, err)
	})

	t.Run("handle credentials.json reading error", func(t *testing.T) {
		service, authImpl := setup()
		given := errors.New("read file error")
		authImpl.ReadFile = func(filename string) (bytes []byte, err error) {
			return nil, given
		}
		err := service.Init()
		assert.Error(t, err)
		assert.EqualError(t, err, given.Error())
	})

	t.Run("handle json config parsing error", func(t *testing.T) {
		service, authImpl := setup()
		given := errors.New("config parse from json")
		authImpl.ConfigFromJSON = func(jsonKey []byte, scope ...string) (config *oauth2.Config, err error) {
			return nil, given
		}
		err := service.Init()
		assert.Error(t, err)
		assert.EqualError(t, err, given.Error())
	})

	t.Run("GetToken error", func(t *testing.T) {
		service, authImpl := setup()
		given := errors.New("get token error")
		authImpl.GetToken = func(config *oauth2.Config, fileName string, f api.Func) (token *oauth2.Token, err error) {
			return nil, given
		}
		err := service.Init()
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

			actual, err := api.GetTokenWithBrowser(ln)
			if tc.success {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.success, actual == want)
		})
	}
}
