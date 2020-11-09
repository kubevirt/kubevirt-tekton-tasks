package vm

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dv"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	cdicliv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
	"time"
)

func WaitForVM(kubevirtClient kubevirtcliv1.KubevirtClient,
	cdiClientSet cdicliv1beta1.CdiV1beta1Interface,
	namespace, name string,
	vmiPhase kubevirtv1.VirtualMachineInstancePhase,
	pvcsAreDataVolumes bool,
	timeout time.Duration) (*kubevirtv1.VirtualMachine, error) {
	var numOfVMPollsBeforeError = 5 / (constants.PollInterval / time.Second) // 5 sec

	var vm *kubevirtv1.VirtualMachine

	err := wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		var err error
		vm, err = kubevirtClient.VirtualMachine(namespace).Get(name, &metav1.GetOptions{})
		if err != nil {
			if numOfVMPollsBeforeError == 0 {
				return true, err
			}
			numOfVMPollsBeforeError--
			return false, nil
		}

		// check DataVolumes' successes
		for _, volume := range vm.Spec.Template.Spec.Volumes {
			var name string
			if dataVolume := volume.DataVolume; dataVolume != nil {
				name = dataVolume.Name
			}
			if pvc := volume.PersistentVolumeClaim; pvcsAreDataVolumes && pvc != nil {
				name = pvc.ClaimName
			}

			if name != "" && !dv.IsDataVolumeImportSuccessful(cdiClientSet, namespace, name) {
				return false, nil
			}
		}

		if vmiPhase != "" {
			vmi, err := kubevirtClient.VirtualMachineInstance(namespace).Get(name, &metav1.GetOptions{})
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
