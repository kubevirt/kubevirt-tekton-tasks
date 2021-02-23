package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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

	if data.PublicKeySecretNamespace == "" {
		data.PublicKeySecretNamespace = options.ResolveNamespace(data.PublicKeySecretTargetNamespace)
	}

	if data.PrivateKeySecretName != "" {
		data.PrivateKeySecretName = E2ETestsRandomName(data.PrivateKeySecretName)
	}

	if data.PrivateKeySecretNamespace == "" {
		data.PrivateKeySecretNamespace = options.ResolveNamespace(data.PrivateKeySecretTargetNamespace)
	}

}

func (c *GenerateSshKeysTestConfig) GetTaskRun() *v1beta1.TaskRun {
	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName("taskrun-generate-ssh-keys"),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: GenerateSshKeysClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: c.ServiceAccount,
			Params: []v1beta1.Param{
				{
					Name: GenerateSshKeysParams.PublicKeySecretName,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.PublicKeySecretName,
					},
				},
				{
					Name: GenerateSshKeysParams.PublicKeySecretNamespace,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.PublicKeySecretNamespace,
					},
				},
				{
					Name: GenerateSshKeysParams.PrivateKeySecretName,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.PrivateKeySecretName,
					},
				},
				{
					Name: GenerateSshKeysParams.PrivateKeySecretNamespace,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.PrivateKeySecretNamespace,
					},
				},
				{
					Name: GenerateSshKeysParams.PrivateKeyConnectionOptions,
					Value: v1beta1.ArrayOrString{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: c.TaskData.PrivateKeyConnectionOptions,
					},
				},
				{
					Name: GenerateSshKeysParams.AdditionalSSHKeygenOptions,
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: c.TaskData.AdditionalSSHKeygenOptions,
					},
				},
			},
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
