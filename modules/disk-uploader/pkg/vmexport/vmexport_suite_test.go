package vmexport_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVMExport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VMExport Suite")
}
