package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zrma/uds-go/api"
	"github.com/zrma/uds-go/mocks"
	"golang.org/x/oauth2"
)

var _ = Describe("파일 읽기 테스트", func() {
	var service *api.Service
	BeforeEach(func() {
		reader := &mocks.Reader{}
		reader.ReadFileReturns([]byte{}, nil)
		reader.ConfigFromJSONReturns(&oauth2.Config{}, nil)

		author := &mocks.Author{}
		author.GetTokenReturns(&oauth2.Token{})

		service = &api.Service{
			Reader: reader,
			Author: author,
		}
	})

	It("test", func() {
		err := service.Init()
		Expect(err).ShouldNot(HaveOccurred())
	})
})
