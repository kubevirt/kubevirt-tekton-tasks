package execattributes_test

import (
	"os"
	"path"
	"reflect"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

type resultChecker func(result interface{})

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

	DescribeTable("Init fails", func(expectedErrMessage string, secretPath string, secretSetup map[string]string, expectedAttributes map[string]interface{}) {
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
			} else if reflect.TypeOf(expectedValue).Kind() == reflect.Func {
				if receive, ok := expectedValue.(resultChecker); ok {
					receive(results[0].Interface())
				} else {
					Fail("invalid expectedValue func")
				}
			} else {
				Expect(results[0].Interface()).To(Equal(expectedValue))
			}
		}
		log.Logger().Debug(CurrentSpecReport().FullText(), zap.Object("execAttributes", attributes)) // test MarshalLogObject
	},
		Entry("secret missing", "secret does not exist", "invalid/path/to/secret", map[string]string{}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		Entry("secret type file missing", "type secret attribute is required", "", map[string]string{}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		Entry("secret type missing", "type secret attribute is required", "", map[string]string{
			"type": "",
		}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		Entry("invalid secret type", "is invalid type", "", map[string]string{
			"type": "authenticate",
		}, map[string]interface{}{
			"GetType":          constants.ExecSecretType(""),
			"GetSSHAttributes": nil,
		}),
		Entry("empty ssh type", "", "", map[string]string{
			"type": "ssh",
		}, map[string]interface{}{
			"GetType":          constants.SSHSecretType,
			"GetSSHAttributes": execattributes.NewSSHAttributes(),
		}),
		Entry("empty ssh type detected via ssh-privatekey", "", "", map[string]string{
			"ssh-privatekey": testconstants.SSHTestPrivateKey,
		}, map[string]interface{}{
			"GetType": constants.SSHSecretType,
			"GetSSHAttributes": resultChecker(func(result interface{}) {
				if sshAttributes, ok := result.(execattributes.SSHAttributes); ok {
					Expect(sshAttributes).ShouldNot(BeNil())
					Expect(sshAttributes.GetPrivateKey()).Should(Equal(testconstants.SSHTestPrivateKey))
				}
			}),
		}),
	)
})
