package api_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zrma/uds-go/api"
	"github.com/zrma/uds-go/mocks"
	"golang.org/x/oauth2"
)

var _ = Describe("파일 읽기 테스트", func() {
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

	It("성공", func() {
		err := service.Init()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("파일 읽기 예외 처리 확인", func() {
		expected := errors.New("read file error")
		author.ReadFileReturns(nil, expected)
		err := service.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err).Should(Equal(expected))
	})

	It("JSON config 예외 처리 확인", func() {
		expected := errors.New("config parse from json")
		author.ConfigFromJSONReturns(nil, expected)
		err := service.Init()
		Expect(err).Should(HaveOccurred())
		Expect(err).Should(Equal(expected))
	})
})
