package testconfigs

import (
	"strings"

	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GenerateSshKeysTaskData struct {
	PublicKeySecretTargetNamespace  TargetNamespace
	PrivateKeySecretTargetNamespace TargetNamespace

	PublicKeySecretName         string
	PublicKeySecretNamespace    string
	PrivateKeySecretName        string
	PrivateKeySecretNamespace   string
	PrivateKeyConnectionOptions []string
	AdditionalSSHKeygenOptions  string
}

type GenerateSshKeysTestConfig struct {
	TaskRunTestConfig
	TaskData GenerateSshKeysTaskData

	deploymentNamespace string
}

func (d *GenerateSshKeysTaskData) GetPrivateKeyConnectionOptions() map[string]string {
	result := make(map[string]string, len(d.PrivateKeyConnectionOptions))

	for _, keyVal := range d.PrivateKeyConnectionOptions {
		if split := strings.SplitN(keyVal, ":", 2); len(split) == 2 {
			key := strings.TrimSpace(split[0])
			result[key] = split[1]
		}
	}
	return result
}

func (c *GenerateSshKeysTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace

	data := &c.TaskData

	if data.PublicKeySecretName != "" {
		data.PublicKeySecretName = E2ETestsRandomName(data.PublicKeySecretName)
	}

	data.PublicKeySecretNamespace = options.GetDeployNamespace()

	if data.PrivateKeySecretName != "" {
		data.PrivateKeySecretName = E2ETestsRandomName(data.PrivateKeySecretName)
	}

	data.PrivateKeySecretNamespace = options.GetDeployNamespace()

}

func (c *GenerateSshKeysTestConfig) GetTaskRun() *pipev1.TaskRun {
	params := []pipev1.Param{
		{
			Name: GenerateSshKeysParams.PublicKeySecretName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.PublicKeySecretName,
			},
		},
		{
			Name: GenerateSshKeysParams.PublicKeySecretNamespace,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.PublicKeySecretNamespace,
			},
		},
		{
			Name: GenerateSshKeysParams.PrivateKeySecretName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.PrivateKeySecretName,
			},
		},
		{
			Name: GenerateSshKeysParams.PrivateKeySecretNamespace,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.PrivateKeySecretNamespace,
			},
		},
		{
			Name: GenerateSshKeysParams.AdditionalSSHKeygenOptions,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: c.TaskData.AdditionalSSHKeygenOptions,
			},
		},
	}

	if len(c.TaskData.PrivateKeyConnectionOptions) > 0 {
		params = append(params, pipev1.Param{
			Name: GenerateSshKeysParams.PrivateKeyConnectionOptions,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: c.TaskData.PrivateKeyConnectionOptions,
			},
		})
	}

	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-generate-ssh-keys"),
			Namespace: c.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: GenerateSshKeysTaskName,
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			Params:  params,
		},
	}
}
func (c *GenerateSshKeysTestConfig) GetExpectedResults() map[string]string {
	expectedResults := map[string]string{
		GenerateSshKeysResults.PublicKeySecretName:       c.TaskData.PublicKeySecretName,
		GenerateSshKeysResults.PublicKeySecretNamespace:  c.TaskData.PublicKeySecretNamespace,
		GenerateSshKeysResults.PrivateKeySecretName:      c.TaskData.PrivateKeySecretName,
		GenerateSshKeysResults.PrivateKeySecretNamespace: c.TaskData.PrivateKeySecretNamespace,
	}

	for _, namespaceKey := range []string{GenerateSshKeysResults.PublicKeySecretNamespace, GenerateSshKeysResults.PrivateKeySecretNamespace} {
		if expectedResults[namespaceKey] == "" {
			expectedResults[namespaceKey] = c.deploymentNamespace
		}
	}

	for key, _ := range expectedResults {
		if expectedResults[key] == "" {
			delete(expectedResults, key)
		}
	}

	return expectedResults
}
