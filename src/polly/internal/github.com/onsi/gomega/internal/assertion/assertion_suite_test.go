package assertion_test

import (
	. "polly/internal/github.com/onsi/ginkgo"
	. "polly/internal/github.com/onsi/gomega"

	"testing"
)

func TestAssertion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Assertion Suite")
}
