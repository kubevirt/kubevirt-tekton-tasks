package execattributes_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path"
	"reflect"
)

const testPublicKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7nkqvSU0JHqtIZTZfhKaiQZeZAPsVg7RwZxSnSnvUxL8hj+G/MpefM8B3ctReJ6lM80JUJO62dKiHvON3QNF14ig4aquGlg15czMPvG8QSyAMxzRZs0dBiLXFvUy+6/6+cNERY3QkbH6pjc3GUON4c/5oRvgXhHHQx+darr6VnZcHGu8h3CwR3IDBnfuhEkeCaIPqZ5RSW1b+znWHRM728gv6A7b4XXUaqwNiF0xkdAdzMwK2/d1hJ9QVRqzfrIcJeK3FRoZfM4QEEv7IlufjdJ0RquRgK+E/bbnHTS+K2T+JrvLyerGgMCv5NOX1Z7GvJKvgLdpeGLW+WTzGnyFtD7ZVn0tnhYQhBAbbChzMYD16yykbxXpn3EV7Bcy4gMB5jVAOgy6t26EuDts27EEQCGbT4NZyQ3chlAI9nTy8bFxH2RznU1u2u3MsBDhe8z0QHGT0u2H6jDWHq8BD5JTRt3gMbjoST0izEXAkn8QjDnAhhJGF3IrpAUQsb+NrTCk= test@test"
const testPrivateKey = `
-----BEGIN OPENSSH PRIVATE KEY-----
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
OKHf1VObOAzTzFjwAAAAlhbnN5QGFuc3k=
-----END OPENSSH PRIVATE KEY-----
`

var _ = Describe("SSHAttributes", func() {
	var testSecretPath string

	BeforeEach(func() {
		testSecretPath = path.Join(testPath, testconstants.TestRandomName("ssh-attr-secret"))
		err := os.MkdirAll(testSecretPath, testDirMode)
		Expect(err).Should(Succeed())
	})

	AfterEach(func() {
		err := os.RemoveAll(testSecretPath)
		Expect(err).Should(Succeed())
	})

	table.DescribeTable("Init fails", func(expectedErrMessage string, secretSetup map[string]string) {
		secretSetup["type"] = "ssh"

		PrepareTestSecret(testSecretPath, secretSetup)
		attributes := execattributes.NewExecAttributes()

		err := attributes.Init(testSecretPath)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
		log.GetLogger().Debug(CurrentGinkgoTestDescription().FullTestText, zap.Object("execAttributes", attributes)) // test MarshalLogObject
	},
		table.Entry("user missing", "user secret attribute is required", map[string]string{}),
		table.Entry("private-key missing", "private-key secret attribute is required", map[string]string{
			"user": "root",
		}),
		table.Entry("public key missing", "host-public-key or disable-strict-host-key-checking=true secret attribute is required", map[string]string{
			"user":        "root",
			"private-key": testPrivateKey,
		}),
		table.Entry("no port argument", "ssh option requires an argument -- p", map[string]string{
			"user":                   "fedora",
			"private-key":            testPrivateKey,
			"host-public-key":        testPublicKey,
			"additional-ssh-options": "-C -p",
		}),
		table.Entry("bad port argument", "Bad port '22.0'", map[string]string{
			"user":                   "fedora",
			"private-key":            testPrivateKey,
			"host-public-key":        testPublicKey,
			"additional-ssh-options": "-C -p 22.0",
		}),
	)

	table.DescribeTable("test various sshAttributes", func(secretSetup map[string]string, expectedAttributes map[string]interface{}) {
		secretSetup["type"] = "ssh"

		PrepareTestSecret(testSecretPath, secretSetup)
		attributes := execattributes.NewExecAttributes()

		err := attributes.Init(testSecretPath)
		Expect(err).Should(Succeed())

		sshAttributes := attributes.GetSSHAttributes()

		for methodName, expectedValue := range expectedAttributes {
			results := reflect.ValueOf(sshAttributes).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}

		log.GetLogger().Info(CurrentGinkgoTestDescription().FullTestText, zap.Object("execAttributes", attributes)) // test MarshalLogObject

	},
		table.Entry("minimal setup", map[string]string{
			"user":            "fedora",
			"private-key":     testPrivateKey,
			"host-public-key": testPublicKey,
		}, map[string]interface{}{
			"GetUser":                      "fedora",
			"GetPort":                      22,
			"GetAdditionalSSHOptions":      "-oStrictHostKeyChecking=yes",
			"GetPrivateKey":                testPrivateKey,
			"GetHostPublicKey":             testPublicKey,
			"GetStrictHostKeyCheckingMode": "yes",
		}),
		table.Entry("disable strict host key checking + custom options", map[string]string{
			"user":                             "fedora",
			"private-key":                      testPrivateKey,
			"disable-strict-host-key-checking": "true",
			"additional-ssh-options":           "-C -p 8022",
		}, map[string]interface{}{
			"GetUser":                      "fedora",
			"GetPort":                      8022,
			"GetAdditionalSSHOptions":      "-C -p 8022 -oStrictHostKeyChecking=accept-new",
			"GetPrivateKey":                testPrivateKey,
			"GetHostPublicKey":             "",
			"GetStrictHostKeyCheckingMode": "accept-new",
		}),
		table.Entry("invalid disable-strict-host-key-checking value", map[string]string{
			"user":                             "fedora",
			"private-key":                      testPrivateKey,
			"host-public-key":                  testPublicKey,
			"disable-strict-host-key-checking": "yes", // should be true
		}, map[string]interface{}{
			"GetUser":                      "fedora",
			"GetPort":                      22,
			"GetAdditionalSSHOptions":      "-oStrictHostKeyChecking=yes",
			"GetPrivateKey":                testPrivateKey,
			"GetHostPublicKey":             testPublicKey,
			"GetStrictHostKeyCheckingMode": "yes",
		}),
	)

	It("does common operations correctly", func() {
		PrepareTestSecret(testSecretPath, map[string]string{
			"type":            "ssh",
			"user":            "fedora",
			"private-key":     testPrivateKey,
			"host-public-key": testPublicKey,
		})
		attributes := execattributes.NewExecAttributes()

		err := attributes.Init(testSecretPath)
		Expect(err).Should(Succeed())

		sshAttributes := attributes.GetSSHAttributes()

		// GetSSHDir
		current, err := user.Current()
		Expect(err).Should(Succeed())
		homeDir := current.HomeDir
		Expect(sshAttributes.GetSSHDir()).To(Equal(path.Join(homeDir, ".ssh")))

		// GetSSHExecutableName
		Expect(sshAttributes.GetSSHExecutableName()).To(Equal("ssh"))

		// IncludesSSHOption
		Expect(sshAttributes.IncludesSSHOption("StrictHostKeyChecking")).To(BeTrue())
		Expect(sshAttributes.IncludesSSHOption("CheckHostIP")).To(BeFalse())

		// IncludesSSHOption and AddAdditionalSSHOption
		sshAttributes.AddAdditionalSSHOption("CheckHostIP", "yes")
		Expect(sshAttributes.IncludesSSHOption("StrictHostKeyChecking")).To(BeTrue())
		Expect(sshAttributes.IncludesSSHOption("CheckHostIP")).To(BeTrue())
		Expect(sshAttributes.GetAdditionalSSHOptions()).Should(Equal("-oStrictHostKeyChecking=yes -oCheckHostIP=yes"))

		log.GetLogger().Info(CurrentGinkgoTestDescription().FullTestText, zap.Object("execAttributes", attributes)) // test MarshalLogObject
	})
})
