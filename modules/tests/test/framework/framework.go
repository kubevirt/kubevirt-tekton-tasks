package framework

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	templatev1 "github.com/openshift/api/template/v1"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/tekton"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

var TestOptionsInstance = &testoptions.TestOptions{}
var ClientsInstance = &Clients{}

type ManagedResources struct {
	taskRun     *pipev1beta1.TaskRun
	dataVolumes []*cdiv1beta1.DataVolume
	vms         []*kubevirtv1.VirtualMachine
	templates   []*templatev1.Template
}
type Framework struct {
	*testoptions.TestOptions
	*Clients

	managedResources ManagedResources
}

type TestConfig interface {
	GetLimitScope() constants.TestScope
	Init(options *testoptions.TestOptions)
}

func NewFramework() *Framework {
	f := &Framework{
		TestOptions: TestOptionsInstance,
		Clients:     ClientsInstance,
	}

	AfterEach(f.AfterEach)
	return f
}

func (f *Framework) TestSetup(config TestConfig) {
	limitScope := config.GetLimitScope()
	if limitScope != "" && limitScope != f.Scope {
		Skip(fmt.Sprintf("runs only in %v scope", limitScope))
	}
	config.Init(f.TestOptions)
}

func (f *Framework) AfterEach() {
	failed := CurrentGinkgoTestDescription().Failed
	taskRun := f.managedResources.taskRun
	hasTaskRun := taskRun != nil

	if failed {
		defer func() {
			if hasTaskRun && !f.Debug {
				defer f.TknClient.TaskRuns(taskRun.Namespace).Delete(taskRun.Name, &metav1.DeleteOptions{})
			}
			tekton.PrintTaskRunDebugInfo(f.TknClient, f.CoreV1Client, taskRun.Namespace, taskRun.Name)
		}()
	}

	if f.Debug {
		// leave resources alive for inspection
		return
	}

	if hasTaskRun && !failed { // failed has its own cleanup
		defer f.TknClient.TaskRuns(taskRun.Namespace).Delete(taskRun.Name, &metav1.DeleteOptions{})
	}

	for _, dv := range f.managedResources.dataVolumes {
		defer f.CdiClient.DataVolumes(dv.Namespace).Delete(dv.Name, &metav1.DeleteOptions{})
	}
	for _, vm := range f.managedResources.vms {
		defer f.KubevirtClient.VirtualMachine(vm.Namespace).Delete(vm.Name, &metav1.DeleteOptions{})
	}
	for _, t := range f.managedResources.templates {
		defer f.TemplateClient.Templates(t.Namespace).Delete(t.Name, &metav1.DeleteOptions{})
	}
}

func (f *Framework) ManageTaskRun(taskRun *pipev1beta1.TaskRun) *Framework {
	f.managedResources.taskRun = taskRun
	return f
}

func (f *Framework) ManageDataVolumes(dataVolumes ...*cdiv1beta1.DataVolume) *Framework {
	for _, dataVolume := range dataVolumes {
		if dataVolume != nil && dataVolume.Name != "" && dataVolume.Namespace != "" {
			f.managedResources.dataVolumes = append(f.managedResources.dataVolumes, dataVolume)
		}
	}
	return f
}

func (f *Framework) ManageVMs(vms ...*kubevirtv1.VirtualMachine) *Framework {
	for _, vm := range vms {
		if vm != nil && vm.Name != "" && vm.Namespace != "" {
			f.managedResources.vms = append(f.managedResources.vms, vm)
		}
	}
	return f
}

func (f *Framework) ManageTemplates(templatest ...*templatev1.Template) *Framework {
	for _, t := range templatest {
		if t != nil && t.Name != "" && t.Namespace != "" {
			f.managedResources.templates = append(f.managedResources.templates, t)
		}
	}
	return f
}
