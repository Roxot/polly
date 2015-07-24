package reporters_test

import (
	. "polly/internal/github.com/onsi/ginkgo"
	. "polly/internal/github.com/onsi/gomega"

	"testing"
)

func TestReporters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reporters Suite")
}
