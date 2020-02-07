package api_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-test/deep"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"

	"github.com/zrma/uds-go/pkg/api"
)

var _ = Describe("GetToken", func() {
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

	Context("GetToken", func() {
		It("should fail", func() {
			author := api.AuthImpl{}
			actual, err := author.GetToken(&config, "", func() (s string, e error) {
				return "", errors.New("error-1234")
			})
			Expect(err).Should(HaveOccurred())
			Expect(actual).Should(BeNil())
		})
	})

	Context("GetTokenFromWeb", func() {
		It("read string failed", func() {
			token, err := api.GetTokenFromWeb(&config, func() (s string, e error) {
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

		err := api.SaveToken(tokenPath, &expected)
		Expect(err).ShouldNot(HaveOccurred())
		actual, err := api.TokenFromFile(tokenPath)
		Expect(err).ShouldNot(HaveOccurred())
		diff := deep.Equal(*actual, expected)
		Expect(diff).Should(BeNil())

		By("GetToken test after setting files...")
		author := api.AuthImpl{}
		actual, err = author.GetToken(&oauth2.Config{}, tokenPath, func() (s string, e error) {
			return
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(actual).ShouldNot(BeNil())
		diff = deep.Equal(*actual, expected)
		Expect(diff).Should(BeNil())
	})
})
