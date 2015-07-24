package codelocation_test

import (
	. "polly/internal/github.com/onsi/ginkgo"
	. "polly/internal/github.com/onsi/gomega"

	"testing"
)

func TestCodelocation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CodeLocation Suite")
}
