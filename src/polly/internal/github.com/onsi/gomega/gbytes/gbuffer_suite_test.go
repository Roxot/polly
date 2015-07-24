package gbytes_test

import (
	. "polly/internal/github.com/onsi/ginkgo"
	. "polly/internal/github.com/onsi/gomega"

	"testing"
)

func TestGbytes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gbytes Suite")
}
