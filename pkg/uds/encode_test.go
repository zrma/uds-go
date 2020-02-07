package uds

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("validate base64 encoding", func() {
	str := `abc123!?$*&()'-=@~`
	base64 := `YWJjMTIzIT8kKiYoKSctPUB+`

	It("should encode", func() {
		Expect(encode([]byte(str))).Should(Equal(base64))
	})
	It("should decode", func() {
		actual, err := decode(base64)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(actual).Should(Equal([]byte(str)))
	})
})
