package templates

import (
	"context"
	"encoding/json"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-vm-template/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	k8sv1 "k8s.io/api/core/v1"

	templatev1 "github.com/openshift/api/template/v1"
	v1 "github.com/openshift/api/template/v1"
	tempclient "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	templateclientset "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

type templateProvider struct {
	client tempclient.TemplateV1Interface
}

type TemplateProvider interface {
	Get(string, string) (*templatev1.Template, error)
	Patch(*v1.Template) (*templatev1.Template, error)
}

func NewTemplateProvider(client tempclient.TemplateV1Interface) TemplateProvider {
	return &templateProvider{
		client: client,
	}
}

func (t *templateProvider) Get(namespace string, name string) (*templatev1.Template, error) {
	return t.client.Templates(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (t *templateProvider) Patch(template *v1.Template) (*templatev1.Template, error) {
	data, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	return t.client.Templates(template.Namespace).Patch(context.TODO(), template.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
}

type TemplateUpdator struct {
	cliOptions       *parse.CLIOptions
	templateProvider TemplateProvider
}

func NewTemplateUpdator(cliOptions *parse.CLIOptions) (*TemplateUpdator, error) {
	log.Logger().Debug("initialized clients and providers")

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return &TemplateUpdator{
		cliOptions:       cliOptions,
		templateProvider: NewTemplateProvider(templateclientset.NewForConfigOrDie(config)),
	}, nil
}

func (t *TemplateUpdator) ModifyTemplate() (*v1.Template, error) {
	log.Logger().Debug("retrieving template", zap.String("name", t.cliOptions.GetTemplateName()), zap.String("namespace", t.cliOptions.GetTemplateNamespace()))
	template, err := t.templateProvider.Get(t.cliOptions.GetTemplateNamespace(), t.cliOptions.GetTemplateName())
	if err != nil {
		return nil, err
	}

	updatedTemplate, err := t.UpdateTemplate(template)
	if err != nil {
		return nil, err
	}

	return t.templateProvider.Patch(updatedTemplate)
}

func (t *TemplateUpdator) UpdateTemplate(template *v1.Template) (*v1.Template, error) {
	t.setValuesToTemplate(template)
	vms, err := DecodeVM(template)
	if err != nil {
		return nil, err
	}
	updatedVM := t.setValuesToVMs(vms)

	return EncodeVMToTemplate(template, updatedVM)
}

func (t *TemplateUpdator) setValuesToTemplate(template *v1.Template) {
	labels := t.cliOptions.GetVMLabels()
	if template.Labels == nil && len(labels) > 0 {
		template.Labels = make(map[string]string)
	}

	for key, value := range labels {
		template.Labels[key] = value
	}

	annotations := t.cliOptions.GetVMAnnotations()
	if template.Annotations == nil && len(annotations) > 0 {
		template.Annotations = make(map[string]string)
	}

	for key, value := range annotations {
		template.Annotations[key] = value
	}
}

func (t *TemplateUpdator) setValuesToVMs(vms []*kubevirtv1.VirtualMachine) []*kubevirtv1.VirtualMachine {
	updatedVMs := make([]*kubevirtv1.VirtualMachine, len(vms))
	labels := t.cliOptions.GetVMLabels()
	annotations := t.cliOptions.GetVMAnnotations()

	for i, vm := range vms {
		if vm.Labels == nil && len(labels) > 0 {
			vm.Labels = make(map[string]string)
		}

		for key, value := range labels {
			vm.Labels[key] = value
		}

		if vm.Annotations == nil && len(annotations) > 0 {
			vm.Annotations = make(map[string]string)
		}

		for key, value := range annotations {
			vm.Annotations[key] = value
		}

		if vm.Spec.Template.Spec.Domain.CPU == nil {
			vm.Spec.Template.Spec.Domain.CPU = &kubevirtv1.CPU{}
		}

		if sockets := t.cliOptions.GetCPUSockets(); sockets > 0 {
			vm.Spec.Template.Spec.Domain.CPU.Sockets = uint32(sockets)
		}
		if cores := t.cliOptions.GetCPUCores(); cores > 0 {
			vm.Spec.Template.Spec.Domain.CPU.Cores = uint32(cores)
		}
		if threads := t.cliOptions.GetCPUThreads(); threads > 0 {
			vm.Spec.Template.Spec.Domain.CPU.Threads = uint32(threads)
		}
		if memory := t.cliOptions.GetMemory(); memory != nil {
			vm.Spec.Template.Spec.Domain.Resources.Requests[k8sv1.ResourceMemory] = *memory
		}
		updatedVMs[i] = vm
	}

	return updatedVMs
}

func EncodeVMToTemplate(template *templatev1.Template, vms []*kubevirtv1.VirtualMachine) (*v1.Template, error) {
	vmsRaw := make([]runtime.RawExtension, len(vms))
	for i, vm := range vms {
		raw, err := json.Marshal(vm)
		if err != nil {
			return nil, err
		}
		vmsRaw[i].Raw = raw
	}
	template.Objects = vmsRaw
	return template, nil
}

func DecodeVM(template *templatev1.Template) ([]*kubevirtv1.VirtualMachine, error) {
	var vms []*kubevirtv1.VirtualMachine

	for _, obj := range template.Objects {
		decoder := kubevirtv1.Codecs.UniversalDecoder(kubevirtv1.GroupVersion)
		decoded, err := runtime.Decode(decoder, obj.Raw)
		if err != nil {
			return nil, err
		}
		vm, ok := decoded.(*kubevirtv1.VirtualMachine)
		if ok {
			vms = append(vms, vm)
			break
		}
	}
	if len(vms) == 0 {
		return nil, zerrors.NewMissingRequiredError("no VM object found in the template")
	}
	return vms, nil
}
