package uds_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUds(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Uds Suite")
}
