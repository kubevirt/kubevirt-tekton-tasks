package testconfigs

import (
	"strconv"

	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	v1 "github.com/openshift/api/template/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ModifyTemplateTaskData struct {
	Template *v1.Template

	TemplateName             string
	SourceTemplateNamespace  TargetNamespace
	CPUCores                 string
	CPUSockets               string
	CPUThreads               string
	Memory                   string
	TemplateNamespace        string
	TemplateLabels           []string
	TemplateAnnotations      []string
	VMAnnotations            []string
	VMLabels                 []string
	Disks                    []string
	Volumes                  []string
	DataVolumeTemplates      []string
	TemplateParameters       []string
	DeleteDatavolumeTemplate bool
	DeleteDisks              bool
	DeleteVolumes            bool
	DeleteTemplateParameters bool
}

type ModifyTemplateTestConfig struct {
	TaskRunTestConfig
	TaskData ModifyTemplateTaskData

	deploymentNamespace string
}

func (m *ModifyTemplateTestConfig) Init(options *testoptions.TestOptions) {
	m.deploymentNamespace = options.DeployNamespace
	m.TaskData.TemplateNamespace = options.ResolveNamespace(m.TaskData.SourceTemplateNamespace, options.TestNamespace)

	if m.TaskData.Template != nil {
		m.TaskData.Template.Name = E2ETestsRandomName(m.TaskData.Template.Name)
		if m.TaskData.TemplateName != "" {
			m.TaskData.TemplateName = m.TaskData.Template.Name
		}
		m.TaskData.Template.Namespace = m.TaskData.TemplateNamespace
	}
}

func (m *ModifyTemplateTestConfig) GetTaskRun() *v1beta1.TaskRun {
	params := []v1beta1.Param{
		{
			Name: TemplateNameOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: m.TaskData.TemplateName,
			},
		}, {
			Name: TemplateNamespaceOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: m.TaskData.TemplateNamespace,
			},
		}, {
			Name: CPUCoresOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: m.TaskData.CPUCores,
			},
		}, {
			Name: CPUSocketsOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: m.TaskData.CPUSockets,
			},
		}, {
			Name: CPUThreadsOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: m.TaskData.CPUThreads,
			},
		}, {
			Name: MemoryOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: m.TaskData.Memory,
			},
		}, {
			Name: TemplateLabelsOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.TemplateLabels,
			},
		}, {
			Name: TemplateAnnotationsOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.TemplateAnnotations,
			},
		}, {
			Name: VMLabelsOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.VMLabels,
			},
		}, {
			Name: VMAnnotationsOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.VMAnnotations,
			},
		}, {
			Name: DisksOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.Disks,
			},
		}, {
			Name: VolumesOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.Volumes,
			},
		}, {
			Name: DataVolumeTemplatesOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.DataVolumeTemplates,
			},
		}, {
			Name: TemplateParametersOptionName,
			Value: v1beta1.ArrayOrString{
				Type:     v1beta1.ParamTypeArray,
				ArrayVal: m.TaskData.TemplateParameters,
			},
		}, {
			Name: DeleteDatavolumeTemplateOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteDatavolumeTemplate),
			},
		}, {
			Name: DeleteDisksOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteDisks),
			},
		}, {
			Name: DeleteVolumesOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteVolumes),
			},
		}, {
			Name: DeleteTemplateParametersOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteTemplateParameters),
			},
		},
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(ModifyTemplateTaskRunName),
			Namespace: m.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: ModifyTemplateClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: m.GetTaskRunTimeout()},
			ServiceAccountName: m.ServiceAccount,
			Params:             params,
		},
	}
}
