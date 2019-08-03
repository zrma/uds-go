package api_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zrma/uds-go/api"
	"github.com/zrma/uds-go/mocks"
	"golang.org/x/oauth2"
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
