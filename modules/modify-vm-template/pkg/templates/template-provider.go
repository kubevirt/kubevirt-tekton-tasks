package templates

import (
	"context"
	"encoding/json"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-vm-template/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
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
	vm, err := zutils.DecodeVM(template)
	if err != nil {
		return nil, err
	}
	updatedVM := t.setValuesToVM(vm)

	return EncodeVMToTemplate(template, updatedVM)
}

func (t *TemplateUpdator) setValuesToTemplate(template *v1.Template) {
	labels := t.cliOptions.GetVMLabels()
	template.Labels = t.concatMaps(template.Labels, labels)

	annotations := t.cliOptions.GetVMAnnotations()
	template.Annotations = t.concatMaps(template.Annotations, annotations)
}

func (t *TemplateUpdator) concatMaps(a, b map[string]string) map[string]string {
	lenB := len(b)
	if a == nil && lenB > 0 {
		a = make(map[string]string, lenB)
	}

	for key, value := range b {
		a[key] = value
	}
	return a
}

func (t *TemplateUpdator) setValuesToVM(vm *kubevirtv1.VirtualMachine) *kubevirtv1.VirtualMachine {
	labels := t.cliOptions.GetVMLabels()
	annotations := t.cliOptions.GetVMAnnotations()

	vm.Labels = t.concatMaps(vm.Labels, labels)

	vm.Annotations = t.concatMaps(vm.Annotations, annotations)

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

	return vm
}

func EncodeVMToTemplate(template *templatev1.Template, vm *kubevirtv1.VirtualMachine) (*v1.Template, error) {
	objectsRaw := make([]runtime.RawExtension, 1)

	raw, err := json.Marshal(vm)
	if err != nil {
		return nil, err
	}
	objectsRaw[0].Raw = raw

	template.Objects = objectsRaw
	return template, nil
}
