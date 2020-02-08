package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"

	"github.com/zrma/uds-go/pkg/api"
)

var _ = Describe("Service", func() {
	var service *api.Service
	var authImpl *api.AuthImpl
	BeforeEach(func() {
		authImpl = &api.AuthImpl{
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
		service = &api.Service{
			AuthImpl: authImpl,
		}
	})

	It("success to initialize", func() {
		err := service.Init()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("handle credentials.json reading error", func() {
		expected := errors.New("read file error")
		authImpl.ReadFile = func(filename string) (bytes []byte, err error) {
			return nil, expected
		}
		err := service.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err).Should(Equal(expected))
	})

	It("handle json config parsing error", func() {
		expected := errors.New("config parse from json")
		authImpl.ConfigFromJSON = func(jsonKey []byte, scope ...string) (config *oauth2.Config, err error) {
			return nil, expected
		}
		err := service.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err).Should(Equal(expected))
	})

	It("GetToken error", func() {
		expected := errors.New("get token error")
		authImpl.GetToken = func(config *oauth2.Config, fileName string, f api.Func) (token *oauth2.Token, err error) {
			return nil, expected
		}
		err := service.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err).Should(Equal(expected))
	})
})
