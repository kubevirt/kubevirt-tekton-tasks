package parse_test

import (
	"reflect"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/parse"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

var (
	defaultNS = "default"
)

var _ = Describe("CLIOptions", func() {
	DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init().Error()).To(ContainSubstring(expectedErrMessage))
	},
		Entry("invalid public secret name", "invalid public-key-secret-name value: a lowercase RFC 1123 subdomain must consist of", &parse.CLIOptions{
			PublicKeySecretName: "invalid name",
		}),
		Entry("invalid public secret namespace", "invalid public-key-secret-namespace value: a lowercase RFC 1123 subdomain must consist of", &parse.CLIOptions{
			PublicKeySecretNamespace: "invalid ns",
		}),
		Entry("invalid private secret name", "invalid private-key-secret-name value: a lowercase RFC 1123 subdomain must consist of", &parse.CLIOptions{
			PrivateKeySecretName: "%invalid-name",
		}),
		Entry("invalid private secret namespace", "invalid private-key-secret-namespace value: a lowercase RFC 1123 subdomain must consist of", &parse.CLIOptions{
			PrivateKeySecretNamespace: "(-invalid-ns",
		}),
		Entry("invalid connection options 1", "invalid private-key connection options: no key found before \"root\"; pair should be in \"KEY:VAL\" format", &parse.CLIOptions{
			PrivateKeyConnectionOptions: []string{"root", "K2=V2"},
		}),
		Entry("invalid connection options 2", "invalid private-key connection options: no key found before \":root\"; pair should be in \"KEY:VAL\" format", &parse.CLIOptions{
			PrivateKeyConnectionOptions: []string{":root"},
		}),
	)

	DescribeTable("Parses and returns correct values", func(options *parse.CLIOptions, expectedOptions map[string]interface{}) {
		Expect(options.Init()).Should(Succeed())

		for methodName, expectedValue := range expectedOptions {
			results := reflect.ValueOf(options).MethodByName(methodName).Call([]reflect.Value{})
			Expect(results[0].Interface()).To(Equal(expectedValue))
		}
	},
		Entry("returns valid defaults", &parse.CLIOptions{
			PublicKeySecretNamespace:  defaultNS,
			PrivateKeySecretNamespace: defaultNS,
		}, map[string]interface{}{
			"GetPublicKeySecretName":         "",
			"GetPublicKeySecretNamespace":    defaultNS,
			"GetPrivateKeySecretName":        "",
			"GetPrivateKeySecretNamespace":   defaultNS,
			"GetSshKeygenOptions":            "",
			"GetPrivateKeyConnectionOptions": map[string]string{},
			"GetDebugLevel":                  zapcore.InfoLevel,
		}),
		Entry("handles cli arguments + trim", &parse.CLIOptions{
			PublicKeySecretName:         "test-public ",
			PublicKeySecretNamespace:    " my-test-ns",
			PrivateKeySecretName:        "test-private",
			PrivateKeySecretNamespace:   "   my-other-ns  ",
			PrivateKeyConnectionOptions: []string{" user:root", "additional-ssh-options:-p 8022"},
			SshKeygenOptions:            "-t ed25519 -a 100 ",
			Debug:                       true,
		}, map[string]interface{}{
			"GetPublicKeySecretName":       "test-public",
			"GetPublicKeySecretNamespace":  "my-test-ns",
			"GetPrivateKeySecretName":      "test-private",
			"GetPrivateKeySecretNamespace": "my-other-ns",
			"GetSshKeygenOptions":          "-t ed25519 -a 100 ",
			"GetPrivateKeyConnectionOptions": map[string]string{
				"user":                   "root",
				"additional-ssh-options": "-p 8022",
			},
			"GetDebugLevel": zapcore.DebugLevel,
		}),
	)
})
