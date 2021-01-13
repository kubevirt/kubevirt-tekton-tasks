package execattributes_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utilstest"
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	testPath     = "./unit-test-tmp/"
	testFileMode = 0666
	testDirMode  = 0777
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExecAttributes Suite")
}

var _ = BeforeSuite(func() {
	utilstest.SetupTestSuite()

	err := os.MkdirAll(testPath, testDirMode)
	Expect(err).Should(Succeed())
})
var _ = AfterSuite(func() {
	utilstest.TearDownSuite()
	err := os.RemoveAll(testPath)
	Expect(err).Should(Succeed())
})

func PrepareTestSecret(basePath string, setup map[string]string) {
	for filename, content := range setup {
		err := writeToFile(path.Join(basePath, filename), content)
		Expect(err).Should(Succeed())
	}
}

func writeToFile(filename string, content string) error {
	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC

	f, err := os.OpenFile(filename, flags, testFileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(content)); err != nil {
		return err
	}
	return f.Sync()
}
