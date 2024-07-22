package test

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

const (
	helloWorldScript = `#!/bin/bash
echo hello world
`
	failScript = `#!/bin/bash
echo fail
exit 25
`
	sleepScript = `#!/bin/bash
sleep 30
echo fail
exit 5
`
)

var _ = Describe("Execute in VM / Cleanup VM", func() {
	f := framework.NewFramework()

	BeforeEach(func() {
		if f.TestOptions.SkipExecuteInVMTests {
			Skip("skipExecuteInVMTests is set to true, skipping tests")
		}
	})

	sshConnectionInfo := map[string]string{
		"type":                             "ssh",
		"user":                             "fedora",
		"ssh-privatekey":                   testconstants.SSHTestPrivateKey,
		"disable-strict-host-key-checking": "true",
	}
	fedoraCloudConfig := testobjects.CloudConfig{
		SSHAuthorizedKeys: []string{
			testconstants.SSHTestPublicKey,
		},
	}
	for _, c := range []ExecInVMMode{ExecuteInVMMode, CleanupVMMode} {
		execInVMMode := c

		DescribeTable(string(execInVMMode), func(config *testconfigs.ExecuteOrCleanupVMTestConfig) {
			config.TaskData.ExecInVMMode = execInVMMode
			f.TestSetup(config)

			if secret := config.TaskData.Secret; secret != nil {
				secret, err := f.K8sClient.CoreV1().Secrets(secret.Namespace).Create(context.Background(), secret, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageSecrets(secret)
			}

			if vm := config.TaskData.VM; vm != nil {
				// put large cloudInits to secrets
				for _, volume := range vm.Spec.Template.Spec.Volumes {
					if volume.CloudInitNoCloud != nil && len([]byte(volume.CloudInitNoCloud.UserData)) > 2048 {
						cloudInitSecret := &corev1.Secret{
							ObjectMeta: metav1.ObjectMeta{
								Name: vm.Name + "-" + volume.Name,
							},
							StringData: map[string]string{
								"userdata": volume.CloudInitNoCloud.UserData,
							},
						}
						cloudInitSecret, err := f.K8sClient.CoreV1().Secrets(vm.Namespace).Create(context.Background(), cloudInitSecret, v1.CreateOptions{})
						Expect(err).ShouldNot(HaveOccurred())
						f.ManageSecrets(cloudInitSecret)
						volume.CloudInitNoCloud.UserData = ""
						volume.CloudInitNoCloud.UserDataSecretRef = &corev1.LocalObjectReference{
							Name: cloudInitSecret.Name,
						}
					}
				}
				vm, err := f.KubevirtClient.VirtualMachine(vm.Namespace).Create(context.Background(), vm, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageVMs(vm)
				if config.TaskData.ShouldStartVM {
					err := f.KubevirtClient.VirtualMachine(vm.Namespace).Start(context.Background(), vm.Name, &kubevirtv1.StartOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					time.Sleep(Timeouts.WaitBeforeExecutingVM.Duration)
				}
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccessOrFailure(config.ExpectSuccess).
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectTermination(config.ExpectedTermination).
				ExpectResults(nil)
		},
			// negative cases
			Entry("no vm", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "missing value for vm-name option",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					Secret: testobjects.NewTestSecret(map[string]string{}),
				},
			}),
			Entry("invalid vm name", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "vm-name is not a valid name: a lowercase RFC 1123 subdomain must consist of",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VMName: "name with spaces",
					Secret: testobjects.NewTestSecret(map[string]string{}),
				},
			}),
			Entry("no action", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "no action was specified: at least one of the following options is required: command|script|stop|delete",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:     testobjects.NewTestFedoraCloudVM("no-action").Build(),
					Secret: testobjects.NewTestSecret(map[string]string{}),
				},
			}),
			Entry("too many actions", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "only one of command|script options is allowed",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:      testobjects.NewTestFedoraCloudVM("too-many-actions").Build(),
					Secret:  testobjects.NewTestSecret(map[string]string{}),
					Script:  helloWorldScript,
					Command: []string{"echo"},
				},
			}),
			Entry("no secret", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "connection secret should not be empty",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:         testobjects.NewTestFedoraCloudVM("no-secret").Build(),
					SecretName: "__empty__",
					Script:     helloWorldScript,
				},
			}),
			Entry("no secret type", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "type secret attribute is required",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:     testobjects.NewTestFedoraCloudVM("no-secret-type").Build(),
					Secret: testobjects.NewTestSecret(map[string]string{}),
					Script: helloWorldScript,
				},
			}),
			Entry("invalid secret type", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "Ssh is invalid type",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("invalid-secret-type").Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type": "Ssh",
					}),
					Script: helloWorldScript,
				},
			}),
			Entry("no secret private-key", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "ssh-privatekey secret attribute is required",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("no-secret-private-key").Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type": "ssh",
					}),
					Script: helloWorldScript,
				},
			}),
			Entry("no secret user", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "user secret attribute is required",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("no-secret-user").Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type":           "ssh",
						"ssh-privatekey": testconstants.SSHTestPrivateKey,
					}),
					Script: helloWorldScript,
				},
			}),
			Entry("no secret host-key", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "host-public-key or disable-strict-host-key-checking=true secret attribute is required",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("no-secret-host-key").Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type":           "ssh",
						"user":           "fedora",
						"ssh-privatekey": testconstants.SSHTestPrivateKey,
					}),
					Script: helloWorldScript,
				},
			}),
			Entry("non existent VM", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "virtualmachine.kubevirt.io \"non-existent\" not found",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					Secret: testobjects.NewTestSecret(sshConnectionInfo),
					VMName: "non-existent",
					Script: helloWorldScript,
				},
			}),
			Entry("not working VM", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:     testobjects.NewTestFedoraCloudVM("not-working-vm").WithMemory("5000Pi").Build(),
					Secret: testobjects.NewTestSecret(sshConnectionInfo),
					Script: helloWorldScript,
				},
			}),
			Entry("not authorized VM", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "Permission denied",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:     testobjects.NewTestFedoraCloudVM("not-authorized-vm").Build(),
					Secret: testobjects.NewTestSecret(sshConnectionInfo),
					Script: helloWorldScript,
				},
			}),
			Entry("fail script", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					Timeout:      Timeouts.QuickTaskRun,
					ExpectedLogs: "fail",
					ExpectedTermination: &testconfigs.TaskRunExpectedTermination{
						ExitCode: 25,
					},
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:     testobjects.NewTestFedoraCloudVM("fail-script").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret: testobjects.NewTestSecret(sshConnectionInfo),
					Script: failScript,
				},
			}),
			Entry("execute with wrong public key", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "REMOTE HOST IDENTIFICATION HAS CHANGED",
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("execute-with-wrong-public-key").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type":            "ssh",
						"user":            "fedora",
						"ssh-privatekey":  testconstants.SSHTestPrivateKey,
						"host-public-key": testconstants.SSHTestPublicKey2,
					}),
					Script: helloWorldScript,
				},
			}),
			// positive cases
			Entry("execute script", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello world",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:     testobjects.NewTestFedoraCloudVM("execute-script").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret: testobjects.NewTestSecret(sshConnectionInfo),
					Script: helloWorldScript,
				},
			}),
			Entry("execute command in running vm", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:            testobjects.NewTestFedoraCloudVM("execute-command-in-running-vm").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret:        testobjects.NewTestSecret(sshConnectionInfo),
					ShouldStartVM: true,
					Command:       []string{"echo", "hello"},
				},
			}),
			Entry("execute command with args", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello world !",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:          testobjects.NewTestFedoraCloudVM("execute-command-with-args").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret:      testobjects.NewTestSecret(sshConnectionInfo),
					Command:     []string{"echo", "hello"},
					CommandArgs: []string{"world", "!"},
				},
			}),
			Entry("execute command in the same namespace as deploy", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello world",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM:     testobjects.NewTestFedoraCloudVM("execute-command-in-same-ns").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret: testobjects.NewTestSecret(sshConnectionInfo),
					Script: helloWorldScript,
				},
			}),
			Entry("execute script with options", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello world",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("execute-script-with-options").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type":           "ssh",
						"user":           "fedora",
						"ssh-privatekey": testconstants.SSHTestPrivateKey,
						// TODO change to safer accept-new once a newer version of ssh which supports this option is available in CI
						"additional-ssh-options":           "-oStrictHostKeyChecking=no -C",
						"disable-strict-host-key-checking": "true",
					}),
					Script: helloWorldScript,
				},
			}),
			Entry("execute with public key", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello world",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("execute-with-public-key").WithCloudConfig(testobjects.CloudConfig{
						SSHAuthorizedKeys: []string{
							testconstants.SSHTestPublicKey,
						},
						SSHKeys: testobjects.CloudConfigSSHKeys{
							RSAPrivate: testconstants.SSHTestPrivateKey2,
							RSAPublic:  testconstants.SSHTestPublicKey2,
						},
					}).Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type":            "ssh",
						"user":            "fedora",
						"ssh-privatekey":  testconstants.SSHTestPrivateKey,
						"host-public-key": testconstants.SSHTestPublicKey2,
					}),
					Script: helloWorldScript,
				},
			}),
			Entry("execute with malformed private key", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello world",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("execute-with-malformed-private-key").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret: testobjects.NewTestSecret(map[string]string{
						"type":                             "ssh",
						"user":                             "fedora",
						"ssh-private-key":                  testconstants.SSHTestPrivateKeyWithoutLastNewLine,
						"disable-strict-host-key-checking": "true",
					}),
					Script: helloWorldScript,
				},
			}),
			Entry("execute with kubernetes.io/ssh-auth secret type", &testconfigs.ExecuteOrCleanupVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs:  "hello world",
					ExpectSuccess: true,
				},
				TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
					VM: testobjects.NewTestFedoraCloudVM("execute-with-kubernetes-ssh-secret-type").WithCloudConfig(fedoraCloudConfig).Build(),
					Secret: &corev1.Secret{
						ObjectMeta: v1.ObjectMeta{
							Name:      "testsecret",
							Namespace: testconstants.NamespaceTestDefault,
						},
						StringData: map[string]string{
							"user":                             "fedora",
							corev1.SSHAuthPrivateKey:           testconstants.SSHTestPrivateKey,
							"disable-strict-host-key-checking": "true",
						},
						Type: corev1.SecretTypeSSHAuth,
					},
					Script: helloWorldScript,
				},
			}),
		)
	}

	DescribeTable("cleanup vm actions", func(config *testconfigs.ExecuteOrCleanupVMTestConfig) {
		config.TaskData.ExecInVMMode = CleanupVMMode
		f.TestSetup(config)

		if secret := config.TaskData.Secret; secret != nil {
			secret, err := f.K8sClient.CoreV1().Secrets(secret.Namespace).Create(context.Background(), secret, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageSecrets(secret)
		}

		if vm := config.TaskData.VM; vm != nil {
			vm, err := f.KubevirtClient.VirtualMachine(vm.Namespace).Create(context.Background(), vm, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageVMs(vm)
			if config.TaskData.ShouldStartVM {
				err := f.KubevirtClient.VirtualMachine(vm.Namespace).Start(context.Background(), vm.Name, &kubevirtv1.StartOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				time.Sleep(Timeouts.WaitBeforeExecutingVM.Duration)
			}
		}

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccessOrFailure(config.ExpectSuccess).
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectTermination(config.ExpectedTermination).
			ExpectResults(nil)

		vm, err := f.KubevirtClient.VirtualMachine(config.TaskData.VMNamespace).Get(context.Background(), config.TaskData.VMName, metav1.GetOptions{})

		if config.TaskData.Delete {
			Expect(err).Should(HaveOccurred())
		} else if config.TaskData.Stop {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*vm.Spec.Running).To(BeFalse())
		}
	},
		// negative cases
		Entry("execute and stops vm with too low timeout", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectedLogs: "command timed out",
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("execute-too-low-timeout-stop-vm").WithCloudConfig(fedoraCloudConfig).Build(),
				Secret:        testobjects.NewTestSecret(sshConnectionInfo),
				Script:        sleepScript,
				ShouldStartVM: true,
				Stop:          true,
				Timeout: &metav1.Duration{
					Duration: 27 * time.Second,
				},
			},
		}),

		Entry("starts and execute and stops vm with too low timeout", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectedLogs: "command timed out",
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:     testobjects.NewTestFedoraCloudVM("start-execute-too-low-timeout-stop-vm").WithCloudConfig(fedoraCloudConfig).Build(),
				Secret: testobjects.NewTestSecret(sshConnectionInfo),
				Script: sleepScript,
				Stop:   true,
				Timeout: &metav1.Duration{
					Duration: 27 * time.Second,
				},
			},
		}),
		// positive cases
		Entry("stop vm", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("stop-vm").Build(),
				SecretName:    "__empty__",
				ShouldStartVM: true,
				Stop:          true,
			},
		}),
		Entry("stop non running vm", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:         testobjects.NewTestFedoraCloudVM("stop-non-running-vm").Build(),
				SecretName: "__empty__",
				Stop:       true,
			},
		}),
		Entry("delete vm", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("delete-vm").Build(),
				SecretName:    "__empty__",
				ShouldStartVM: true,
				Delete:        true,
			},
		}),
		Entry("stop and delete vm", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("stop-delete-vm").Build(),
				SecretName:    "__empty__",
				ShouldStartVM: true,
				Stop:          true,
				Delete:        true,
			},
		}),
		Entry("execute and stop and delete vm", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectedLogs:  "hello world",
				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("execute-stop-delete-vm").WithCloudConfig(fedoraCloudConfig).Build(),
				Secret:        testobjects.NewTestSecret(sshConnectionInfo),
				Script:        helloWorldScript,
				ShouldStartVM: true,
				Stop:          true,
				Delete:        true,
			},
		}),
		Entry("execute and stops vm with timeout", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectedLogs:  "hello world",
				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("execute-timeout-stop-vm").WithCloudConfig(fedoraCloudConfig).Build(),
				Secret:        testobjects.NewTestSecret(sshConnectionInfo),
				Script:        helloWorldScript,
				ShouldStartVM: true,
				Stop:          true,
				Timeout:       Timeouts.WaitForVMStart,
			},
		}),
		Entry("execute and deletes vm with timeout", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectedLogs:  "hello world",
				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("execute-timeout-delete-vm").WithCloudConfig(fedoraCloudConfig).Build(),
				Secret:        testobjects.NewTestSecret(sshConnectionInfo),
				Script:        helloWorldScript,
				ShouldStartVM: true,
				Delete:        true,
				Timeout:       Timeouts.WaitForVMStart,
			},
		}),
		Entry("stops failed VMI", &testconfigs.ExecuteOrCleanupVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{

				ExpectSuccess: true,
			},
			TaskData: testconfigs.ExecuteOrCleanupVMTaskData{
				VM:            testobjects.NewTestFedoraCloudVM("stops-failed-vmi").WithMemory("100Pi").Build(),
				SecretName:    "__empty__",
				ShouldStartVM: true,
				Stop:          true,
			},
		}),
	)
})
