package api_test

import (
	"errors"
	"fmt"
	"github.com/go-test/deep"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zrma/uds-go/api"
	"github.com/zrma/uds-go/mocks"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"os"
	"path/filepath"
	"time"
)

var _ = Describe("Service", func() {
	var service *api.Service
	author := &mocks.Author{}
	BeforeEach(func() {
		author.ReadFileReturns([]byte{}, nil)
		author.ConfigFromJSONReturns(&oauth2.Config{}, nil)
		author.GetTokenReturns(&oauth2.Token{})

		service = &api.Service{
			Author: author,
		}
	})

	It("success to initialize", func() {
		err := service.Init()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("handle credentials.json reading error", func() {
		expected := errors.New("read file error")
		author.ReadFileReturns(nil, expected)
		err := service.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err).Should(Equal(expected))
	})

	It("handle json config parsing error", func() {
		expected := errors.New("config parse from json")
		author.ConfigFromJSONReturns(nil, expected)
		err := service.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err).Should(Equal(expected))
	})
})

var _ = Describe("GetTokenFromWeb", func() {
	It("read string failed", func() {
		token, err := api.GetTokenFromWeb(&oauth2.Config{
			ClientID:     "client-1",
			ClientSecret: "secret-2",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "auth-url-1",
				TokenURL:  "token-url-2",
				AuthStyle: 0,
			},
			RedirectURL: "localhost-3",
			Scopes:      []string{drive.DriveScope},
		}, func() (s string, e error) {
			return "", errors.New("test")
		})
		Expect(err).Should(HaveOccurred())
		Expect(token).Should(BeNil())
	})

	It("exchange failed", func() {
		token, err := api.GetTokenFromWeb(&oauth2.Config{
			ClientID:     "client-1",
			ClientSecret: "secret-2",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "auth-url-1",
				TokenURL:  "token-url-2",
				AuthStyle: 0,
			},
			RedirectURL: "localhost-3",
			Scopes:      []string{drive.DriveScope},
		}, func() (s string, e error) {
			return "token-1234", nil
		})
		Expect(err).Should(HaveOccurred())
		Expect(token).Should(BeNil())
	})
})

var _ = Describe("token file I/O", func() {
	const (
		prefix    = "tmp_"
		tokenName = "token1234.json"
	)
	tmpPath := os.TempDir()
	tmpPath = filepath.Join(tmpPath, "uds-go")
	tokenPath := filepath.Join(tmpPath, prefix+tokenName)

	BeforeEach(func() {
		if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
			err := os.Mkdir(tmpPath, os.ModePerm)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("create", tmpPath)
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpPath)
		Expect(err).ShouldNot(HaveOccurred())

		fmt.Println("remove", tmpPath)
	})

	It("token file save/load", func() {
		expected := oauth2.Token{
			AccessToken:  "token1234",
			TokenType:    "type123",
			RefreshToken: "refresh123",
			Expiry:       time.Now(),
		}

		api.SaveToken(tokenPath, &expected)
		actual, err := api.TokenFromFile(tokenPath)
		Expect(err).ShouldNot(HaveOccurred())
		diff := deep.Equal(*actual, expected)
		Expect(diff).Should(BeNil())

		By("GetToken test after setting files...")
		author := api.AuthorImpl{}
		actual = author.GetToken(&oauth2.Config{}, tokenPath)
		diff = deep.Equal(*actual, expected)
		Expect(diff).Should(BeNil())
	})
})
