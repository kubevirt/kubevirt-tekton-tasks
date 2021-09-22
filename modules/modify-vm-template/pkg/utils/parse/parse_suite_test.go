package parse_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestParse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parse Suite")
}
