package requirements_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testSelector = "metadata.name in (fedora, ubuntu), spec.running != true, invalid.path notin (1, 2, 3), metadata"

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Requirements Suite")
}
