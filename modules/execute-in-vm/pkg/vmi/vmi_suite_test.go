package vmi_test

import (
	"testing"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utilstest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VMI Suite")
}

var _ = BeforeSuite(utilstest.SetupTestSuite)
var _ = AfterSuite(utilstest.TearDownSuite)
