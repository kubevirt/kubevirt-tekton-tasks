package test

import (
	"context"
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Generate SSH Keys", func() {
	f := framework.NewFramework()

	BeforeEach(func() {
		if f.TestOptions.SkipGenerateSSHKeysTests {
			Skip("skipGenerateSSHKeysTests is set to true, skipping tests")
		}
	})

	DescribeTable("taskrun fails and no Secrets are created", func(config *testconfigs.GenerateSshKeysTestConfig) {
		f.TestSetup(config)

		runner := runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			WaitForTaskRunFinish()

		f.ManageSecrets(asManagedSecrets(runner.GetResults())...)

		runner.ExpectFailure().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(nil)
	},
		Entry("invalid public secret name", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   "invalid public-key-secret-name value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PublicKeySecretName: "public secret",
			},
		}),
		Entry("invalid public secret namespace", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   "invalid public-key-secret-namespace value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PublicKeySecretNamespace: "my ns",
			},
		}),
		Entry("invalid private secret name", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   "invalid private-key-secret-name value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretName: "private secret",
			},
		}),
		Entry("invalid private secret namespace", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   "invalid private-key-secret-namespace value: a lowercase RFC 1123 subdomain must consist of",
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretNamespace: "my ns",
			},
		}),
		Entry("invalid ssh-keygen options", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   "unknown option -- 3",
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				AdditionalSSHKeygenOptions: "-3 unknown-param test",
			},
		}),
		Entry("cannot create public secret in different namespace", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountNameNamespaced,
				ExpectedLogs:   "secrets is forbidden",
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PublicKeySecretTargetNamespace: SystemTargetNS,
			},
		}),
		Entry("cannot create private secret in different namespace", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountNameNamespaced,
				ExpectedLogs:   "secrets is forbidden",
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretTargetNamespace: SystemTargetNS,
			},
		}),
	)

	DescribeTable("Secrets are created successfully", func(config *testconfigs.GenerateSshKeysTestConfig) {
		f.TestSetup(config)

		runner := runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			WaitForTaskRunFinish()

		results := runner.GetResults()
		f.ManageSecrets(asManagedSecrets(results)...)

		runner.ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResultsWithLen(config.GetExpectedResults(), 4)

		publicSecret, err := f.CoreV1Client.Secrets(results[GenerateSshKeysResults.PublicKeySecretNamespace]).Get(context.Background(), results[GenerateSshKeysResults.PublicKeySecretName], metav1.GetOptions{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(publicSecret.Data).To(HaveLen(1))

		connectionOptions := config.TaskData.GetPrivateKeyConnectionOptions()

		for _, value := range publicSecret.Data {
			Expect(len(string(value))).Should(BeNumerically(">", 30))
			if connectionOptions[PrivateKeyConnectionOptions.User] != "" {
				Expect(string(value)).To(ContainSubstring(" %v@", connectionOptions[PrivateKeyConnectionOptions.User]))
			}
		}

		privateSecret, err := f.CoreV1Client.Secrets(results[GenerateSshKeysResults.PrivateKeySecretNamespace]).Get(context.Background(), results[GenerateSshKeysResults.PrivateKeySecretName], metav1.GetOptions{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(privateSecret.Type).Should(Equal(corev1.SecretTypeSSHAuth))
		Expect(string(privateSecret.Data[PrivateKeyConnectionOptions.PrivateKey])).Should(ContainSubstring("PRIVATE KEY-----"))
		Expect(len(string(privateSecret.Data[PrivateKeyConnectionOptions.PrivateKey]))).Should(BeNumerically(">", 200))
		Expect(string(privateSecret.Data[PrivateKeyConnectionOptions.Type])).Should(BeEmpty())
		Expect(string(privateSecret.Data[PrivateKeyConnectionOptions.User])).Should(Equal(connectionOptions[PrivateKeyConnectionOptions.User]))
		Expect(string(privateSecret.Data[PrivateKeyConnectionOptions.HostPublicKey])).Should(Equal(connectionOptions[PrivateKeyConnectionOptions.HostPublicKey]))
		Expect(string(privateSecret.Data[PrivateKeyConnectionOptions.DisableStrictHostKeyCheckingAttr])).Should(Equal(connectionOptions[PrivateKeyConnectionOptions.DisableStrictHostKeyCheckingAttr]))
		Expect(string(privateSecret.Data[PrivateKeyConnectionOptions.AdditionalSSHOptionsAttr])).Should(Equal(connectionOptions[PrivateKeyConnectionOptions.AdditionalSSHOptionsAttr]))
	},
		Entry("generates secret with no additional data", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{},
		}),
		Entry("generates secret with private secret name", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretName: "test-private-secret",
			},
		}),
		Entry("generates secret with public secret name", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PublicKeySecretName: "test-public-secret",
			},
		}),
		Entry("generates secret with empty namespaces", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretTargetNamespace: EmptyTargetNS,
				PublicKeySecretTargetNamespace:  DeployTargetNS,
			},
		}),
		Entry("generates complex secret", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretName:           "complex-private-key",
				PublicKeySecretName:            "complex-public-key",
				PublicKeySecretTargetNamespace: DeployTargetNS,
				PrivateKeyConnectionOptions: []string{
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.User, "root"),
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.DisableStrictHostKeyCheckingAttr, "false"),
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.HostPublicKey, testconstants.SSHTestPublicKey),
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.AdditionalSSHOptionsAttr, "-C -p 8022"),
				},
				AdditionalSSHKeygenOptions: "-N \"my long passphrase 4515645\" -C root@myclient",
			},
		}),
		Entry("generates complex mixed secret", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretName:           "complex-private-key",
				PublicKeySecretTargetNamespace: EmptyTargetNS,
				PrivateKeyConnectionOptions: []string{
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.User, "root"),
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.DisableStrictHostKeyCheckingAttr, "true"),
				},
			},
		}),
		Entry("works also in the same namespace as deploy", &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PrivateKeySecretTargetNamespace: DeployTargetNS,
				PublicKeySecretTargetNamespace:  DeployTargetNS,
			},
		}),
	)

	It("appends to existing secret", func() {
		config := &testconfigs.GenerateSshKeysTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: GenerateSshKeysServiceAccountName,
				ExpectedLogs:   ExpectedGenerateSshKeysMessage,
			},
			TaskData: testconfigs.GenerateSshKeysTaskData{
				PublicKeySecretName:            "authorized-keys-append-secret",
				PublicKeySecretTargetNamespace: TestTargetNS,
				PrivateKeyConnectionOptions: []string{
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.User, "root"),
					fmt.Sprintf("%v:%v", PrivateKeyConnectionOptions.DisableStrictHostKeyCheckingAttr, "true"),
				},
			},
		}
		f.TestSetup(config)

		publicSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: config.TaskData.PublicKeySecretName,
			},
			StringData: map[string]string{
				"id_rsa.pub": testconstants.SSHTestPublicKey,
			},
		}

		publicSecret, err := f.CoreV1Client.Secrets(f.TestNamespace).Create(context.Background(), publicSecret, metav1.CreateOptions{})
		Expect(err).ShouldNot(HaveOccurred())
		f.ManageSecrets(publicSecret)

		runner := runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			WaitForTaskRunFinish()

		results := runner.GetResults()
		f.ManageSecrets(asManagedSecrets(results)...)

		runner.ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResultsWithLen(config.GetExpectedResults(), 4)

		publicSecret, err = f.CoreV1Client.Secrets(results[GenerateSshKeysResults.PublicKeySecretNamespace]).Get(context.Background(), results[GenerateSshKeysResults.PublicKeySecretName], metav1.GetOptions{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(publicSecret.Data).To(HaveLen(2))

		// old public key not changed
		Expect(string(publicSecret.Data["id_rsa.pub"])).To(Equal(testconstants.SSHTestPublicKey))

		for _, value := range publicSecret.Data {
			Expect(len(string(value))).Should(BeNumerically(">", 30))
			Expect(string(value)).To(HavePrefix("ssh-rsa"))
			Expect(string(value)).To(ContainSubstring("@"))
		}

		privateSecret, err := f.CoreV1Client.Secrets(results[GenerateSshKeysResults.PrivateKeySecretNamespace]).Get(context.Background(), results[GenerateSshKeysResults.PrivateKeySecretName], metav1.GetOptions{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(privateSecret.Data[PrivateKeyConnectionOptions.PrivateKey])).Should(ContainSubstring("PRIVATE KEY-----"))
	})
})

func asManagedSecrets(results map[string]string) []*corev1.Secret {
	var managedSecrets []*corev1.Secret

	if len(results) > 0 {
		if name, namespace := results[GenerateSshKeysResults.PublicKeySecretName], results[GenerateSshKeysResults.PublicKeySecretNamespace]; name != "" && namespace != "" {
			managedSecrets = append(managedSecrets, &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
			})
		}

		if name, namespace := results[GenerateSshKeysResults.PrivateKeySecretName], results[GenerateSshKeysResults.PrivateKeySecretNamespace]; name != "" && namespace != "" {
			managedSecrets = append(managedSecrets, &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
			})
		}
	}

	return managedSecrets
}
