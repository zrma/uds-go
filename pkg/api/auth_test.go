package api_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"

	"github.com/zrma/uds-go/pkg/api"
)

func TestGetToken(t *testing.T) {
	config := oauth2.Config{
		ClientID:     "client-1",
		ClientSecret: "secret-2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "auth-url-1",
			TokenURL:  "token-url-2",
			AuthStyle: 0,
		},
		RedirectURL: "localhost-3",
		Scopes:      []string{drive.DriveScope},
	}

	t.Run("GetToken should fail", func(t *testing.T) {
		authImpl := api.AuthImpl{
			GetToken: api.GetToken,
		}
		got, err := authImpl.GetToken(&config, "", api.Func{
			TokenFromFile: api.TokenFromFile,
			TokenFromWeb:  api.TokenFromWeb,
			GetAuthCode: func() (s string, e error) {
				return "", errors.New("error-1234")
			},
			SaveToken: api.SaveToken,
		})

		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("TokenFromWeb", func(t *testing.T) {
		t.Run("reading string failed", func(t *testing.T) {
			got, err := api.TokenFromWeb(&config, api.Func{
				GetAuthCode: func() (s string, e error) {
					return "", errors.New("test")
				},
			})
			assert.Error(t, err)
			assert.Nil(t, got)
		})

		t.Run("exchange failed", func(t *testing.T) {
			got, err := api.TokenFromWeb(&oauth2.Config{
				ClientID:     "client-1",
				ClientSecret: "secret-2",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "auth-url-1",
					TokenURL:  "token-url-2",
					AuthStyle: 0,
				},
				RedirectURL: "localhost-3",
				Scopes:      []string{drive.DriveScope},
			}, api.Func{
				GetAuthCode: func() (s string, e error) {
					return "token-1234", nil
				},
			})
			assert.Error(t, err)
			assert.Nil(t, got)
		})
	})
}

func TestTokenFileIO(t *testing.T) {
	tokenPath := setupTmpFolder(t)

	t.Run("token file save/load", func(t *testing.T) {
		want := &oauth2.Token{
			AccessToken:  "token1234",
			TokenType:    "type123",
			RefreshToken: "refresh123",
			Expiry:       time.Now().Add(time.Minute),
		}

		t.Run("save succeed", func(t *testing.T) {
			err := api.SaveToken(tokenPath, want)
			assert.NoError(t, err)

			got, err := api.TokenFromFile(tokenPath)
			assert.NoError(t, err)

			diff := deep.Equal(want, got)
			assert.Nil(t, diff)
		})

		t.Run("GetToken test after setting files...", func(t *testing.T) {
			author := api.AuthImpl{
				GetToken: api.GetToken,
			}
			got, err := author.GetToken(&oauth2.Config{}, tokenPath, api.Func{
				TokenFromFile: api.TokenFromFile,
				TokenFromWeb: func(config *oauth2.Config, f api.Func) (token *oauth2.Token, err error) {
					return want, nil
				},
				GetAuthCode: func() (s string, err error) {
					return "", nil
				},
				SaveToken: api.SaveToken,
			})
			assert.NoError(t, err)
			assert.NotNil(t, got)

			diff := deep.Equal(want, got)
			assert.Nil(t, diff)
		})

		t.Run("token expiration has occurred", func(t *testing.T) {
			want.Expiry = time.Now().Add(-time.Minute)
			err := api.SaveToken(tokenPath, want)
			assert.NoError(t, err)

			got, err := api.TokenFromFile(tokenPath)
			assert.NoError(t, err)

			assert.False(t, got.Valid(), "token has been expired")

			author := api.AuthImpl{
				GetToken: api.GetToken,
			}
			got, err = author.GetToken(&oauth2.Config{}, tokenPath, api.Func{
				TokenFromFile: api.TokenFromFile,
				TokenFromWeb: func(config *oauth2.Config, f api.Func) (token *oauth2.Token, err error) {
					return want, nil
				},
				GetAuthCode: func() (s string, err error) {
					return "", nil
				},
				SaveToken: api.SaveToken,
				FakeExchange: func() (token *oauth2.Token, err error) {
					want.Expiry = time.Now().Add(time.Minute)
					return want, nil
				},
			})

			assert.NoError(t, err)
			assert.True(t, got.Valid(), "token has been refreshed")
		})
	})
}

func setupTmpFolder(t *testing.T) string {
	const (
		prefix    = "tmp_"
		tokenName = "token1234.json"
	)
	tmpPath := os.TempDir()
	tmpPath = filepath.Join(tmpPath, "uds-go")
	tokenPath := filepath.Join(tmpPath, prefix+tokenName)

	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		err := os.Mkdir(tmpPath, os.ModePerm)
		assert.NoError(t, err)

		fmt.Println("create", tmpPath)
	}

	t.Cleanup(func() {
		err := os.RemoveAll(tmpPath)
		assert.NoError(t, err)

		fmt.Println("remove", tmpPath)
	})

	return tokenPath
}
