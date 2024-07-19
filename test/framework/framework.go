package framework

import (
	"context"
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/clients"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/testoptions"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/tekton"
	. "github.com/onsi/ginkgo/v2"
	templatev1 "github.com/openshift/api/template/v1"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	instancetype "kubevirt.io/api/instancetype/v1beta1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var TestOptionsInstance = &testoptions.TestOptions{}
var ClientsInstance = &clients.Clients{}

type ManagedResources struct {
	taskRuns             []*pipev1.TaskRun
	pipelineRuns         []*pipev1.PipelineRun
	pipelines            []*pipev1.Pipeline
	dataVolumes          []*cdiv1beta1.DataVolume
	dataSources          []*cdiv1beta1.DataSource
	vms                  []*kubevirtv1.VirtualMachine
	templates            []*templatev1.Template
	secrets              []*corev1.Secret
	clusterInstancetypes []*instancetype.VirtualMachineClusterInstancetype
}

type Framework struct {
	*testoptions.TestOptions
	*clients.Clients

	managedResources  ManagedResources
	limitEnvScope     constants.EnvScope
	onBeforeTestSetup func(config TestConfig)
}

type TestConfig interface {
	GetLimitEnvScope() constants.EnvScope
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

func (f *Framework) LimitEnvScope(limitEnvScope constants.EnvScope) *Framework {
	if f.limitEnvScope != "" {
		Fail("limitEnvScope was already set")
	}
	f.limitEnvScope = limitEnvScope

	return f
}

func (f *Framework) OnBeforeTestSetup(callback func(config TestConfig)) *Framework {
	f.onBeforeTestSetup = callback
	return f
}

func (f *Framework) TestSetup(config TestConfig) {
	limitEnvScope := config.GetLimitEnvScope()

	// check global env limit first
	if f.limitEnvScope != "" && f.limitEnvScope != f.EnvScope {
		Skip(fmt.Sprintf("runs only in %v", f.limitEnvScope))
	}

	// check test case env limit
	if limitEnvScope != "" && limitEnvScope != f.EnvScope {
		Skip(fmt.Sprintf("runs only in %v", limitEnvScope))
	}

	if f.onBeforeTestSetup != nil {
		f.onBeforeTestSetup(config)
	}
	config.Init(f.TestOptions)
}

func (f *Framework) AfterEach() {
	failed := CurrentSpecReport().Failed()
	taskRuns := f.managedResources.taskRuns
	pipelineRuns := f.managedResources.pipelineRuns

	if failed {
		defer func() {
			if !f.Debug {
				for _, taskRun := range taskRuns {
					defer f.TknClient.TaskRuns(taskRun.Namespace).Delete(context.Background(), taskRun.Name, metav1.DeleteOptions{})
				}
				for _, pipelineRun := range pipelineRuns {
					defer f.TknClient.PipelineRuns(pipelineRun.Namespace).Delete(context.Background(), pipelineRun.Name, metav1.DeleteOptions{})
				}
			}
			for _, taskRun := range taskRuns {
				tekton.PrintTaskRunDebugInfo(f.Clients, taskRun.Namespace, taskRun.Name)
			}
			for _, pipelineRun := range pipelineRuns {
				tekton.PrintPipelineRunDebugInfo(f.Clients, pipelineRun.Namespace, pipelineRun.Name)
			}
		}()
	}

	if f.Debug {
		// leave resources alive for inspection
		return
	}

	if !failed { // failed has its own cleanup
		for _, taskRun := range taskRuns {
			defer f.TknClient.TaskRuns(taskRun.Namespace).Delete(context.Background(), taskRun.Name, metav1.DeleteOptions{})
		}
		for _, pipelineRun := range pipelineRuns {
			defer f.TknClient.PipelineRuns(pipelineRun.Namespace).Delete(context.Background(), pipelineRun.Name, metav1.DeleteOptions{})
		}
	}
	for _, pipeline := range f.managedResources.pipelines {
		defer f.TknClient.Pipelines(pipeline.Namespace).Delete(context.Background(), pipeline.Name, metav1.DeleteOptions{})
	}
	for _, dv := range f.managedResources.dataVolumes {
		defer f.CdiClient.DataVolumes(dv.Namespace).Delete(context.Background(), dv.Name, metav1.DeleteOptions{})
	}
	for _, ds := range f.managedResources.dataSources {
		defer f.CdiClient.DataSources(ds.Namespace).Delete(context.Background(), ds.Name, metav1.DeleteOptions{})
	}
	for _, vm := range f.managedResources.vms {
		defer f.KubevirtClient.VirtualMachine(vm.Namespace).Delete(context.Background(), vm.Name, metav1.DeleteOptions{})
	}
	for _, t := range f.managedResources.templates {
		defer f.TemplateClient.Templates(t.Namespace).Delete(context.Background(), t.Name, metav1.DeleteOptions{})
	}
	for _, s := range f.managedResources.secrets {
		defer f.KubevirtClient.CoreV1().Secrets(s.Namespace).Delete(context.Background(), s.Name, metav1.DeleteOptions{})
	}
	for _, clusterInstancetype := range f.managedResources.clusterInstancetypes {
		defer f.KubevirtClient.VirtualMachineClusterInstancetype().Delete(context.Background(), clusterInstancetype.Name, metav1.DeleteOptions{})
	}
}

func (f *Framework) ManageTaskRuns(taskRuns ...*pipev1.TaskRun) *Framework {
	f.managedResources.taskRuns = append(f.managedResources.taskRuns, taskRuns...)
	return f
}

func (f *Framework) ManagePipelineRuns(pipelineRuns ...*pipev1.PipelineRun) *Framework {
	f.managedResources.pipelineRuns = append(f.managedResources.pipelineRuns, pipelineRuns...)
	return f
}

func (f *Framework) ManagePipelines(pipelines ...*pipev1.Pipeline) *Framework {
	f.managedResources.pipelines = append(f.managedResources.pipelines, pipelines...)
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

func (f *Framework) ManageDataSources(dataSources ...*cdiv1beta1.DataSource) *Framework {
	for _, dataSource := range dataSources {
		if dataSource != nil && dataSource.Name != "" && dataSource.Namespace != "" {
			f.managedResources.dataSources = append(f.managedResources.dataSources, dataSource)
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

func (f *Framework) ManageTemplates(templates ...*templatev1.Template) *Framework {
	for _, t := range templates {
		if t != nil && t.Name != "" && t.Namespace != "" {
			f.managedResources.templates = append(f.managedResources.templates, t)
		}
	}
	return f
}

func (f *Framework) ManageSecrets(secrets ...*corev1.Secret) *Framework {
	for _, s := range secrets {
		if s != nil && s.Name != "" && s.Namespace != "" {
			f.managedResources.secrets = append(f.managedResources.secrets, s)
		}
	}
	return f
}

func (f *Framework) ManageClusterInstancetypes(clusterInstancetypes ...*instancetype.VirtualMachineClusterInstancetype) *Framework {
	for _, clusterInstancetype := range clusterInstancetypes {
		if clusterInstancetype != nil && clusterInstancetype.Name != "" {
			f.managedResources.clusterInstancetypes = append(f.managedResources.clusterInstancetypes, clusterInstancetype)
		}
	}
	return f
}
