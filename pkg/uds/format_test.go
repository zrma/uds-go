package uds

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("validate util function", func() {
	Context("format function", func() {
		type testData struct {
			numOfBytes int64
			expected   string
			ok         bool
		}

		const (
			kb = 1024
			mb = kb * 1024
			gb = mb * 1024
			tb = gb * 1024
		)

		DescribeTable("size table", func(data testData) {
			size, err := format(data.numOfBytes)
			if data.ok {
				Expect(err).ShouldNot(HaveOccurred())
			} else {
				Expect(err).Should(HaveOccurred())
			}
			Expect(size).Should(Equal(data.expected))
		},
			Entry("invalid - negative size", testData{-1, "", false}),
			Entry("", testData{1000, "1000.0 bytes", true}),
			Entry("", testData{kb, "1.0 KB", true}),
			Entry("", testData{800 * kb, "800.0 KB", true}),
			Entry("", testData{mb, "1.0 MB", true}),
			Entry("", testData{tb, "1.0 TB", true}),
			Entry("", testData{1024 * tb, "1024.0 TB", true}),
		)
	})
})
