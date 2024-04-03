package testconfigs

import (
	"strconv"

	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	v1 "github.com/openshift/api/template/v1"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
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
	DeleteTemplate           bool
}

type ModifyTemplateTestConfig struct {
	TaskRunTestConfig
	TaskData ModifyTemplateTaskData

	deploymentNamespace string
}

func (m *ModifyTemplateTestConfig) Init(options *testoptions.TestOptions) {
	m.deploymentNamespace = options.DeployNamespace
	m.TaskData.TemplateNamespace = options.GetDeployNamespace()

	if m.TaskData.Template != nil {
		m.TaskData.Template.Name = E2ETestsRandomName(m.TaskData.Template.Name)
		if m.TaskData.TemplateName != "" {
			m.TaskData.TemplateName = m.TaskData.Template.Name
		}
		m.TaskData.Template.Namespace = m.TaskData.TemplateNamespace
	}
}

func (m *ModifyTemplateTestConfig) GetTaskRun() *pipev1.TaskRun {
	params := []pipev1.Param{
		{
			Name: TemplateNameOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: m.TaskData.TemplateName,
			},
		}, {
			Name: TemplateNamespaceOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: m.TaskData.TemplateNamespace,
			},
		}, {
			Name: CPUCoresOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: m.TaskData.CPUCores,
			},
		}, {
			Name: CPUSocketsOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: m.TaskData.CPUSockets,
			},
		}, {
			Name: CPUThreadsOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: m.TaskData.CPUThreads,
			},
		}, {
			Name: MemoryOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: m.TaskData.Memory,
			},
		}, {
			Name: DeleteDatavolumeTemplateOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteDatavolumeTemplate),
			},
		}, {
			Name: DeleteDisksOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteDisks),
			},
		}, {
			Name: DeleteVolumesOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteVolumes),
			},
		}, {
			Name: DeleteTemplateParametersOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteTemplateParameters),
			},
		}, {
			Name: DeleteTemplateOptionName,
			Value: pipev1.ParamValue{
				Type:      pipev1.ParamTypeString,
				StringVal: strconv.FormatBool(m.TaskData.DeleteTemplate),
			},
		},
	}
	if len(m.TaskData.TemplateLabels) > 0 {
		params = append(params, pipev1.Param{
			Name: TemplateLabelsOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.TemplateLabels,
			},
		})
	}

	if len(m.TaskData.TemplateAnnotations) > 0 {
		params = append(params, pipev1.Param{
			Name: TemplateAnnotationsOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.TemplateAnnotations,
			},
		})
	}

	if len(m.TaskData.VMLabels) > 0 {
		params = append(params, pipev1.Param{
			Name: VMLabelsOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.VMLabels,
			},
		})
	}

	if len(m.TaskData.VMAnnotations) > 0 {
		params = append(params, pipev1.Param{
			Name: VMAnnotationsOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.VMAnnotations,
			},
		})
	}

	if len(m.TaskData.Disks) > 0 {
		params = append(params, pipev1.Param{
			Name: DisksOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.Disks,
			},
		})
	}

	if len(m.TaskData.Volumes) > 0 {
		params = append(params, pipev1.Param{
			Name: VolumesOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.Volumes,
			},
		})
	}

	if len(m.TaskData.DataVolumeTemplates) > 0 {
		params = append(params, pipev1.Param{
			Name: DataVolumeTemplatesOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.DataVolumeTemplates,
			},
		})
	}

	if len(m.TaskData.TemplateParameters) > 0 {
		params = append(params, pipev1.Param{
			Name: TemplateParametersOptionName,
			Value: pipev1.ParamValue{
				Type:     pipev1.ParamTypeArray,
				ArrayVal: m.TaskData.TemplateParameters,
			},
		})
	}

	return &pipev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(ModifyTemplateTaskRunName),
			Namespace: m.deploymentNamespace,
		},
		Spec: pipev1.TaskRunSpec{
			TaskRef: &pipev1.TaskRef{
				Name: ModifyTemplateTaskName,
				Kind: pipev1.NamespacedTaskKind,
			},
			Timeout: &metav1.Duration{Duration: m.GetTaskRunTimeout()},
			Params:  params,
		},
	}
}
