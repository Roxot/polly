package apns_test

import (
	. "polly/internal/github.com/onsi/ginkgo"
	. "polly/internal/github.com/onsi/gomega"

	"testing"
)

func TestApns(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apns Suite")
}
