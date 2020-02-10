package api_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

var _ = Describe("oauth callback server test", func() {
	const expected = "expected"

	type testData struct {
		query   string
		success bool
		err     error
	}

	DescribeTable("test case", func(data testData) {
		ln, err := net.Listen("tcp", ":0")
		Expect(err).ShouldNot(HaveOccurred())

		addr := strings.Split(ln.Addr().String(), "]")[1]
		go func() {
			time.Sleep(500 * time.Millisecond)
			reqUrl := fmt.Sprintf("http://localhost%s/auth/callback/%s", addr, data.query)
			req, err2 := http.NewRequest("GET", reqUrl, nil)
			Expect(err2).ShouldNot(HaveOccurred())

			client := &http.Client{}
			resp, err2 := client.Do(req)
			Expect(err2).ShouldNot(HaveOccurred())

			defer resp.Body.Close()
		}()

		actual, err := api.GetTokenWithBrowser(ln)
		if data.success {
			Expect(err).ShouldNot(HaveOccurred())
		} else {
			Expect(err).Should(HaveOccurred())
		}
		Expect(actual == expected).Should(Equal(data.success))
	},
		Entry("success", testData{
			query:   fmt.Sprintf("?code=%s", expected),
			success: true,
			err:     nil,
		}),
		Entry("empty param", testData{
			query:   "",
			success: false,
			err:     errors.New("invalid callback params \n"),
		}),
	)
})
