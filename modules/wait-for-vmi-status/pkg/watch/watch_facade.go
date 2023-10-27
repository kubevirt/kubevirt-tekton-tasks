package watch

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/requirements"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/parse"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	kubevirtv1 "kubevirt.io/api/core/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	"time"
)

type WatchFacade struct {
	clioptions     *parse.CLIOptions
	kubeClient     *kubernetes.Clientset
	kubevirtClient kubevirtcliv1.KubevirtClient
}

func NewWatchFacade(clioptions *parse.CLIOptions) (*WatchFacade, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kubeClient := kubernetes.NewForConfigOrDie(config)

	kubevirtClient, err := kubevirtcliv1.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubevirt client: %v", err.Error())
	}

	return &WatchFacade{clioptions: clioptions, kubeClient: kubeClient, kubevirtClient: kubevirtClient}, nil
}

func (f *WatchFacade) WaitForVMIConditions() bool {
	successRequirements := f.clioptions.GetSuccessRequirements()
	failureRequirements := f.clioptions.GetFailureRequirements()

	if len(successRequirements) == 0 && len(failureRequirements) == 0 {
		return true
	}

	listerWatcher := cache.NewListWatchFromClient(f.kubevirtClient.RestClient(),
		"virtualmachineinstances",
		f.clioptions.VirtualMachineInstanceNamespace,
		fields.OneTermEqualSelector(v1.ObjectNameField, f.clioptions.VirtualMachineInstanceName),
	)

	stop := make(chan struct{})
	success := make(chan bool, 1)

	eventHandler := func(obj interface{}) {
		if len(successRequirements) > 0 {
			log.Logger().Debug("evaluating condition", zap.String("successCondition", f.clioptions.GetSuccessCondition()))
			if requirements.MatchesRequirements(obj, successRequirements) {
				success <- true
				close(stop)
				return
			}
		}

		if len(failureRequirements) > 0 {
			log.Logger().Debug("evaluating condition", zap.String("failureCondition", f.clioptions.GetFailureCondition()))
			if requirements.MatchesRequirements(obj, failureRequirements) {
				success <- false
				close(stop)
			}
		}
	}

	_, controller := cache.NewInformer(listerWatcher, &kubevirtv1.VirtualMachineInstance{}, time.Second*0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Logger().Debug("vmi added", zap.Reflect("vmi", obj))
			eventHandler(obj)
		},
		DeleteFunc: func(obj interface{}) {
			log.Logger().Debug("vmi deleted", zap.Reflect("vmi", obj))
			eventHandler(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Logger().Debug("vmi changed", zap.Reflect("vmi", newObj))
			eventHandler(newObj)
		},
	})

	controller.Run(stop)

	return <-success
}
