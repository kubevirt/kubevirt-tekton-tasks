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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

type templateProvider struct {
	client tempclient.TemplateV1Interface
}

type TemplateProvider interface {
	Get(string, string) (*templatev1.Template, error)
	Patch(*v1.Template) (*templatev1.Template, error)
	Delete(*v1.Template) error
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

func (t *templateProvider) Delete(template *v1.Template) error {
	return t.client.Templates(template.Namespace).Delete(context.TODO(), template.Name, metav1.DeleteOptions{})
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

	if t.cliOptions.GetDeleteTemplate() {
		if err := t.templateProvider.Delete(template); err != nil {
			return nil, err
		}
		return template, nil
	}

	updatedTemplate, err := t.UpdateTemplate(template)
	if err != nil {
		return nil, err
	}

	return t.templateProvider.Patch(updatedTemplate)
}

func (t *TemplateUpdator) UpdateTemplate(template *v1.Template) (*v1.Template, error) {
	t.setValuesToTemplate(template)
	vm, vmIndex, err := zutils.DecodeVM(template)
	if err != nil {
		return nil, err
	}
	updatedVM := t.setValuesToVM(vm)

	return EncodeVMToTemplate(template, updatedVM, vmIndex)
}

func (t *TemplateUpdator) setValuesToTemplate(template *v1.Template) {
	labels := t.cliOptions.GetTemplateLabels()
	template.Labels = appendToMap(template.Labels, labels)

	annotations := t.cliOptions.GetTemplateAnnotations()
	template.Annotations = appendToMap(template.Annotations, annotations)

	if t.cliOptions.GetDeleteTemplateParameters() {
		template.Parameters = []v1.Parameter{}
	}

	for _, parameter := range t.cliOptions.GetTemplateParameters() {
		replaced := false
		for i, templateParam := range template.Parameters {
			if parameter.Name == templateParam.Name {
				template.Parameters[i] = parameter
				replaced = true
			}
		}
		if !replaced {
			template.Parameters = append(template.Parameters, parameter)
		}
	}
}

func appendToMap(a, b map[string]string) map[string]string {
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

	vm.Labels = appendToMap(vm.Labels, labels)

	vm.Annotations = appendToMap(vm.Annotations, annotations)

	if vm.Spec.Template.Spec.Domain.CPU == nil {
		vm.Spec.Template.Spec.Domain.CPU = &kubevirtv1.CPU{}
	}

	if sockets := t.cliOptions.GetCPUSockets(); sockets > 0 {
		vm.Spec.Template.Spec.Domain.CPU.Sockets = sockets
	}

	if cores := t.cliOptions.GetCPUCores(); cores > 0 {
		vm.Spec.Template.Spec.Domain.CPU.Cores = cores
	}

	if threads := t.cliOptions.GetCPUThreads(); threads > 0 {
		vm.Spec.Template.Spec.Domain.CPU.Threads = threads
	}

	if memory := t.cliOptions.GetMemory(); memory != nil {
		vm.Spec.Template.Spec.Domain.Resources.Requests[k8sv1.ResourceMemory] = *memory
	}

	if t.cliOptions.GetDeleteDisks() {
		vm.Spec.Template.Spec.Domain.Devices.Disks = []kubevirtv1.Disk{}
	}

	if t.cliOptions.GetDeleteVolumes() {
		vm.Spec.Template.Spec.Volumes = []kubevirtv1.Volume{}
	}

	if t.cliOptions.GetDeleteDatavolumeTemplate() {
		deleteDatavolumeTemplateFromVM(vm)
	}

	for _, datavolumeTemplate := range t.cliOptions.GetDatavolumeTemplates() {
		replaced := false
		for i, vmDatavolumeTemplate := range vm.Spec.DataVolumeTemplates {
			if datavolumeTemplate.Name == vmDatavolumeTemplate.Name {
				vm.Spec.DataVolumeTemplates[i] = datavolumeTemplate
				replaced = true
			}
		}
		if !replaced {
			vm.Spec.DataVolumeTemplates = append(vm.Spec.DataVolumeTemplates, datavolumeTemplate)
		}
	}

	for _, disk := range t.cliOptions.GetDisks() {
		replaced := false
		for i, vmDisk := range vm.Spec.Template.Spec.Domain.Devices.Disks {
			if disk.Name == vmDisk.Name {
				vm.Spec.Template.Spec.Domain.Devices.Disks[i] = disk
				replaced = true
			}
		}
		if !replaced {
			vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, disk)
		}
	}

	for _, volume := range t.cliOptions.GetVolumes() {
		replaced := false
		for i, vmVolume := range vm.Spec.Template.Spec.Volumes {
			if volume.Name == vmVolume.Name {
				vm.Spec.Template.Spec.Volumes[i] = volume
				replaced = true
			}
		}
		if !replaced {
			vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volume)
		}
	}

	return vm
}

func EncodeVMToTemplate(template *templatev1.Template, vm *kubevirtv1.VirtualMachine, vmIndex int) (*v1.Template, error) {
	raw, err := json.Marshal(vm)
	if err != nil {
		return nil, err
	}

	template.Objects[vmIndex].Raw = raw
	return template, nil
}

func deleteDatavolumeTemplateFromVM(vm *kubevirtv1.VirtualMachine) {
	dvsToDelete := make(map[string]bool)
	for _, dvTemplate := range vm.Spec.DataVolumeTemplates {
		dvsToDelete[dvTemplate.Name] = true
	}

	if vm.Spec.Template != nil {
		disksToDelete := make(map[string]bool)
		newVolumes := []kubevirtv1.Volume{}
		for _, volume := range vm.Spec.Template.Spec.Volumes {
			if volume.DataVolume != nil {
				if val, ok := dvsToDelete[volume.DataVolume.Name]; ok && val {
					disksToDelete[volume.Name] = true
					continue
				}
			}
			newVolumes = append(newVolumes, volume)
		}
		vm.Spec.Template.Spec.Volumes = newVolumes

		newDisks := []kubevirtv1.Disk{}
		for _, disk := range vm.Spec.Template.Spec.Domain.Devices.Disks {
			if _, ok := disksToDelete[disk.Name]; !ok {
				newDisks = append(newDisks, disk)
			}
		}
		vm.Spec.Template.Spec.Domain.Devices.Disks = newDisks
	}

	vm.Spec.DataVolumeTemplates = []kubevirtv1.DataVolumeTemplateSpec{}
}
