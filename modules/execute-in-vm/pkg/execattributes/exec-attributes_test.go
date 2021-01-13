package execattributes_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"os"
	"path"
	"reflect"
)

var _ = Describe("ExecAttributes", func() {
	var testSecretPath string

	BeforeEach(func() {
		testSecretPath = path.Join(testPath, testconstants.TestRandomName("exec-attr-secret"))
		err := os.MkdirAll(testSecretPath, testDirMode)
		Expect(err).Should(Succeed())
	})

	AfterEach(func() {
		err := os.RemoveAll(testSecretPath)
		Expect(err).Should(Succeed())
	})

	table.DescribeTable("Init fails", func(expectedErrMessage string, secretPath string, secretSetup map[string]string, expectedAttributes map[string]interface{}) {
		if secretPath == "" {
			secretPath = testSecretPath
		}

		PrepareTestSecret(secretPath, secretSetup)
		attributes := execattributes.NewExecAttributes()

		err := attributes.Init(secretPath)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(expectedErrMessage))

		for methodName, expectedValue := range expectedAttributes {
			results := reflect.ValueOf(attributes).MethodByName(methodName).Call([]reflect.Value{})
			if expectedValue == nil {
				Expect(results[0].Interface()).To(BeNil())
			} else {
				Expect(results[0].Interface()).To(Equal(expectedValue))
			}
		}
		log.GetLogger().Debug(CurrentGinkgoTestDescription().FullTestText, zap.Object("execAttributes", attributes)) // test MarshalLogObject
	},
		table.Entry("secret missing", "secret does not exist", "invalid/path/to/secret", map[string]string{}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		table.Entry("secret type file missing", "type secret attribute is required", "", map[string]string{}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		table.Entry("secret type missing", "type secret attribute is required", "", map[string]string{
			"type": "",
		}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		table.Entry("invalid secret type", "is invalid type", "", map[string]string{
			"type": "authenticate",
		}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		table.Entry("empty ssh type", "", "", map[string]string{
			"type": "ssh",
		}, map[string]interface{}{
			"GetType":          constants.SSHSecretType,
			"GetSSHAttributes": execattributes.NewSSHAttributes(),
		}),
	)
})
