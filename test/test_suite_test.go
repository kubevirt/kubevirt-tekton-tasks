package test

import (
	"testing"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/clients"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	. "github.com/onsi/ginkgo/v2"
	ginkgo_reporters "github.com/onsi/ginkgo/v2/reporters"
	v1reporter "kubevirt.io/client-go/reporter"
	qe_reporters "kubevirt.io/qe-tools/pkg/ginkgo-reporters"

	. "github.com/onsi/gomega"
)

var afterSuiteReporters []Reporter

var _ = ReportAfterSuite("Tests", func(report Report) {
	for _, reporter := range afterSuiteReporters {
		ginkgo_reporters.ReportViaDeprecatedReporter(reporter, report)
	}
})

func TestExit(t *testing.T) {
	RegisterFailHandler(Fail)
	BuildTestSuite()

	if qe_reporters.JunitOutput != "" {
		afterSuiteReporters = append(afterSuiteReporters, v1reporter.NewV1JUnitReporter(qe_reporters.JunitOutput))
	}

	if qe_reporters.Polarion.Run {
		afterSuiteReporters = append(afterSuiteReporters, &qe_reporters.Polarion)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Functional test suite")
}

func BuildTestSuite() {
	BeforeSuite(func() {
		err := testoptions.InitTestOptions(framework.TestOptionsInstance)
		noErr(err)
		err = clients.InitClients(framework.ClientsInstance, framework.TestOptionsInstance)
		noErr(err)
	})
}

func noErr(err error) {
	if err != nil {
		Fail(err.Error())
	}
}
