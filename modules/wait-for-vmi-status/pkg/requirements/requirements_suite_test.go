package requirements_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testSelector = "metadata.name in (fedora, ubuntu), spec.runStrategy != true, invalid.path notin (1, 2, 3), metadata"

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Requirements Suite")
}
