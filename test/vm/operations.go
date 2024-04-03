package vm

import (
	"context"
	"fmt"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/dataobject"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	kubevirtv1 "kubevirt.io/api/core/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
)

func WaitForVM(kubevirtClient kubevirtcliv1.KubevirtClient,
	namespace, name string,
	vmiPhase kubevirtv1.VirtualMachineInstancePhase,
	timeout time.Duration,
	skipStorage bool) (*kubevirtv1.VirtualMachine, error) {
	var numOfVMPollsBeforeError = 5 / (constants.PollInterval / time.Second) // 5 sec

	var vm *kubevirtv1.VirtualMachine

	err := wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		var err error
		vm, err = kubevirtClient.VirtualMachine(namespace).Get(context.Background(), name, &metav1.GetOptions{})
		if err != nil {
			if numOfVMPollsBeforeError == 0 {
				return true, err
			}
			numOfVMPollsBeforeError--
			return false, nil
		}

		if !skipStorage {
			// check DataVolumes' successes
			for _, volume := range vm.Spec.Template.Spec.Volumes {
				var name string
				if dataVolume := volume.DataVolume; dataVolume != nil {
					name = dataVolume.Name
				}

				result, err := dataobject.IsDataVolumeImportSuccessful(kubevirtClient, namespace, name)
				if err != nil {
					fmt.Println("error while waiting for datavolume import: ", err.Error())
				}

				if name != "" && !result {
					return false, nil
				}
			}
		}

		if vmiPhase != "" {
			vmi, err := kubevirtClient.VirtualMachineInstance(namespace).Get(context.Background(), name, &metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			return vmi.Status.Phase == vmiPhase, nil
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return vm, nil
}
