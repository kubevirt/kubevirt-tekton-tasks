package fileoptions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFileoptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fileoptions Suite")
}
