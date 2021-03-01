package execattributes_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path"
	"reflect"
)

var _ = Describe("SSHAttributes", func() {
	var testSecretPath string

	BeforeEach(func() {
		testSecretPath = path.Join(testPath, TestRandomName("ssh-attr-secret"))
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
		table.Entry("privatekey missing", "ssh-privatekey secret attribute is required", map[string]string{}),
		table.Entry("user missing", "user secret attribute is required", map[string]string{
			"ssh-privatekey": SSHTestPrivateKey,
		}),
		table.Entry("public key missing", "host-public-key or disable-strict-host-key-checking=true secret attribute is required", map[string]string{
			"user":           "root",
			"ssh-privatekey": SSHTestPrivateKey,
		}),
		table.Entry("no port argument", "ssh option requires an argument -- p", map[string]string{
			"user":                   "fedora",
			"ssh-privatekey":         SSHTestPrivateKey,
			"host-public-key":        SSHTestPublicKey,
			"additional-ssh-options": "-C -p",
		}),
		table.Entry("bad port argument", "Bad port '22.0'", map[string]string{
			"user":                   "fedora",
			"ssh-privatekey":         SSHTestPrivateKey,
			"host-public-key":        SSHTestPublicKey,
			"additional-ssh-options": "-C -p 22.0",
		}),
	)

	table.DescribeTable("test various sshAttributes", func(secretSetup map[string]string, expectedAttributes map[string]interface{}) {
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
			"type":            "ssh",
			"user":            "fedora",
			"ssh-privatekey":  SSHTestPrivateKey,
			"host-public-key": SSHTestPublicKey,
		}, map[string]interface{}{
			"GetUser":                      "fedora",
			"GetPort":                      22,
			"GetAdditionalSSHOptions":      "-oStrictHostKeyChecking=yes",
			"GetPrivateKey":                SSHTestPrivateKey,
			"GetHostPublicKey":             SSHTestPublicKey,
			"GetStrictHostKeyCheckingMode": "yes",
		}),
		table.Entry("minimal setup with alternative private key format", map[string]string{
			"user":            "fedora",
			"ssh-private-key": SSHTestPrivateKey,
			"host-public-key": SSHTestPublicKey,
		}, map[string]interface{}{
			"GetUser":                      "fedora",
			"GetPort":                      22,
			"GetAdditionalSSHOptions":      "-oStrictHostKeyChecking=yes",
			"GetPrivateKey":                SSHTestPrivateKey,
			"GetHostPublicKey":             SSHTestPublicKey,
			"GetStrictHostKeyCheckingMode": "yes",
		}),
		table.Entry("disable strict host key checking + custom options", map[string]string{
			"type":                             "ssh",
			"user":                             "fedora",
			"ssh-privatekey":                   SSHTestPrivateKey,
			"disable-strict-host-key-checking": "true",
			"additional-ssh-options":           "-C -p 8022",
		}, map[string]interface{}{
			"GetUser": "fedora",
			"GetPort": 8022,
			// TODO change to safer acceptNew once a newer version of ssh which supports this option is available in CI
			"GetAdditionalSSHOptions": "-C -p 8022 -oStrictHostKeyChecking=no",
			"GetPrivateKey":           SSHTestPrivateKey,
			"GetHostPublicKey":        "",
			// TODO same here
			"GetStrictHostKeyCheckingMode": "no",
		}),
		table.Entry("invalid disable-strict-host-key-checking value", map[string]string{
			"type":                             "ssh",
			"user":                             "fedora",
			"ssh-privatekey":                   SSHTestPrivateKey,
			"host-public-key":                  SSHTestPublicKey,
			"disable-strict-host-key-checking": "yes", // should be true
		}, map[string]interface{}{
			"GetUser":                      "fedora",
			"GetPort":                      22,
			"GetAdditionalSSHOptions":      "-oStrictHostKeyChecking=yes",
			"GetPrivateKey":                SSHTestPrivateKey,
			"GetHostPublicKey":             SSHTestPublicKey,
			"GetStrictHostKeyCheckingMode": "yes",
		}),
		table.Entry("end newline in private key", map[string]string{
			"type":            "ssh",
			"user":            "fedora",
			"ssh-privatekey":  SSHTestPrivateKeyWithoutLastNewLine,
			"host-public-key": SSHTestPublicKey,
		}, map[string]interface{}{
			"GetPrivateKey": SSHTestPrivateKey,
		}),
	)

	It("does common operations correctly", func() {
		PrepareTestSecret(testSecretPath, map[string]string{
			"type":            "ssh",
			"user":            "fedora",
			"ssh-privatekey":  SSHTestPrivateKey,
			"host-public-key": SSHTestPublicKey,
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
