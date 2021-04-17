package api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

func TestGetToken(t *testing.T) {
	helperBackup := Helper
	t.Cleanup(func() {
		Helper = helperBackup
	})

	Helper.OpenBrowser = nil

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
		Helper.ScanAuthCode = func() (string, error) {
			return "", errors.New("error-1234")
		}

		got, err := Helper.GetToken(&config, "")

		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("getTokenFromWeb", func(t *testing.T) {
		t.Run("reading string failed", func(t *testing.T) {
			Helper.ScanAuthCode = func() (string, error) {
				return "", errors.New("test")
			}

			got, err := Helper.GetTokenFromWeb(&config)
			assert.Error(t, err)
			assert.Nil(t, got)
		})

		t.Run("exchange failed", func(t *testing.T) {
			Helper.ScanAuthCode = func() (string, error) {
				return "token-1234", nil
			}

			got, err := Helper.GetTokenFromWeb(&oauth2.Config{
				ClientID:     "client-1",
				ClientSecret: "secret-2",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "auth-url-1",
					TokenURL:  "token-url-2",
					AuthStyle: 0,
				},
				RedirectURL: "localhost-3",
				Scopes:      []string{drive.DriveScope},
			})
			assert.Error(t, err)
			assert.Nil(t, got)
		})
	})
}

func TestTokenFileIO(t *testing.T) {
	tokenPath := setupTmpFolder(t)

	want := &oauth2.Token{
		AccessToken:  "token1234",
		TokenType:    "type123",
		RefreshToken: "refresh123",
		Expiry:       time.Now().Add(time.Minute),
	}

	t.Run("save succeed", func(t *testing.T) {
		err := saveToken(tokenPath, want)
		assert.NoError(t, err)

		got, err := Helper.GetTokenFromFile(tokenPath)
		assert.NoError(t, err)

		diff := deep.Equal(want, got)
		assert.Nil(t, diff)
	})

	t.Run("GetToken test after setting files...", func(t *testing.T) {
		Helper.ScanAuthCode = func() (string, error) {
			return "", nil
		}
		Helper.GetTokenFromWeb = func(config *oauth2.Config) (token *oauth2.Token, err error) {
			return want, nil
		}

		got, err := Helper.GetToken(&oauth2.Config{}, tokenPath)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		diff := deep.Equal(want, got)
		assert.Nil(t, diff)
	})

	t.Run("token expiration has occurred", func(t *testing.T) {
		givenCfg := &oauth2.Config{}
		Helper.ExchangeToken = func(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
			assert.Equal(t, givenCfg, config)

			equalTokens(t, want, token)

			token.Expiry = time.Now().Add(time.Minute)
			return token, nil
		}
		Helper.ScanAuthCode = func() (string, error) {
			return "", nil
		}
		Helper.GetTokenFromWeb = func(config *oauth2.Config) (*oauth2.Token, error) {
			assert.Equal(t, givenCfg, config)
			return want, nil
		}

		want.Expiry = time.Now().Add(-time.Minute)
		assert.False(t, want.Valid())

		err := saveToken(tokenPath, want)
		assert.NoError(t, err)

		got, err := Helper.GetTokenFromFile(tokenPath)
		assert.NoError(t, err)

		equalTokens(t, want, got)

		assert.False(t, got.Valid(), "token has been expired")

		got, err = Helper.GetToken(givenCfg, tokenPath)

		assert.NoError(t, err)
		assert.True(t, got.Valid(), "token has been refreshed")
	})
}

func equalTokens(t *testing.T, given, got *oauth2.Token) {
	assert.Equal(t, given.TokenType, got.TokenType)
	assert.Equal(t, given.AccessToken, got.AccessToken)
	assert.Equal(t, given.RefreshToken, got.RefreshToken)
	assert.True(t, given.Expiry.Equal(got.Expiry))
}

func setupTmpFolder(t *testing.T) string {
	const (
		tokenName = "token.json"
	)

	fsBackup := AppFs

	AppFs = afero.NewMemMapFs()
	afs := &afero.Afero{Fs: AppFs}

	tmpPath, err := afs.TempDir("tmp", "")
	assert.NoError(t, err)
	tokenPath := filepath.Join(tmpPath, tokenName)

	if _, err := afs.Stat(tmpPath); os.IsNotExist(err) {
		err := afs.Mkdir(tmpPath, os.ModePerm)
		assert.NoError(t, err)

		fmt.Println("create", tmpPath)
	}

	helperBackup := Helper

	t.Cleanup(func() {
		AppFs = fsBackup
		Helper = helperBackup
	})

	return tokenPath
}
