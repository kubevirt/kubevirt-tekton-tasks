package execute

import (
	"context"
	"fmt"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/vmi"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
)

type Executor struct {
	clioptions     *parse.CLIOptions
	kubevirtClient kubecli.KubevirtClient
	executor       RemoteExecutor

	attemptedStart  bool
	attemptedStop   bool
	attemptedDelete bool
	ipAddress       string
}

func NewExecutor(clioptions *parse.CLIOptions, connectionSecretPath string) (*Executor, error) {
	var executor RemoteExecutor

	executor = newSSHExecutor(clioptions, execattributes.NewExecAttributes())
	if clioptions.GetScript() != "" {
		execAttributes := execattributes.NewExecAttributes()

		if err := execAttributes.Init(connectionSecretPath); err != nil {
			return nil, err
		}
		log.Logger().Debug("retrieved connection secret exec attributes", zap.Object("execAttributes", execAttributes))

		switch execAttributes.GetType() {
		case constants.SSHSecretType:
			executor = newSSHExecutor(clioptions, execAttributes)
		default:
			return nil, fmt.Errorf("invalid secret/execution type %v", execAttributes.GetType())
		}
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kubevirtClient, err := kubecli.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", "cannot create kubevirt client", err.Error())
	}

	return &Executor{clioptions: clioptions, kubevirtClient: kubevirtClient, executor: executor}, nil
}

func (e *Executor) EnsureVMRunning(timeout time.Duration) error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()
	logFields := []zap.Field{zap.String("name", vmName), zap.String("namespace", vmNamespace)}

	conditionFn := func() (done bool, err error) {
		vmInstance, err := e.kubevirtClient.VirtualMachineInstance(vmNamespace).Get(context.TODO(), vmName, &v1.GetOptions{})

		if err != nil {
			log.Logger().Debug("could not obtain a vm instance", logFields[0], logFields[1], zap.Error(err))
			switch t := err.(type) {
			case *errors.StatusError:
				if t.Status().Reason == v1.StatusReasonNotFound {
					if err := e.ensureVMStarted(); err != nil {
						return false, err
					}
					log.Logger().Debug(" waiting for a VMI to start", logFields...)
					return false, nil
				}
				return false, err
			default:
				return false, err
			}
		}

		switch vmInstance.Status.Phase {
		case kubevirtv1.Failed:
			log.Logger().Debug("vm instance failed", logFields[0], logFields[1], zap.Reflect("status", vmInstance.Status))
			// maybe the vm just stopped so let's try to start it again
			if err := e.ensureVMStarted(); err != nil {
				return false, err
			}
			log.Logger().Debug("waiting for a VMI to recover", logFields...)
			return false, nil
		case kubevirtv1.Running:
			ipAddress, ipError := vmi.GetPodIPAddress(vmInstance)

			if ipAddress == "" || ipError != nil {
				log.Logger().Debug("ip address not found", logFields[0], logFields[1], zap.Reflect("status", vmInstance.Status))

				if ipError != nil {
					return false, zerrors.NewMissingRequiredError(ipError.Error())
				}
				// wait for ipAddress
				return false, nil
			}
			log.Logger().Debug("ip address found", zap.String("ipAddress", ipAddress))
			e.ipAddress = ipAddress

			return true, nil

		default:
			log.Logger().Debug("waiting for a VMI to start", logFields...)
			return false, nil
		}
	}

	if timeout <= 0 {
		return wait.PollImmediateInfinite(constants.PollVMIInterval, conditionFn)
	}
	return wait.PollImmediate(constants.PollVMIInterval, timeout, conditionFn)

}

func (e *Executor) EnsureVMStopped() error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()

	return wait.PollImmediateInfinite(constants.PollVMItoStopInterval, func() (bool, error) {
		vmi, err := e.kubevirtClient.VirtualMachineInstance(vmNamespace).Get(context.TODO(), vmName, &v1.GetOptions{})

		if err == nil {
			switch vmi.Status.Phase {
			case kubevirtv1.Succeeded, kubevirtv1.Failed:
				return true, nil
			}

			if stopErr := e.ensureVMStop(); stopErr != nil {
				switch t := stopErr.(type) {
				case *errors.StatusError:
					if t.Status().Reason == v1.StatusReasonConflict { // stop already requested
						return false, nil
					}
					return false, stopErr
				default:
					return false, stopErr
				}
			}

			log.Logger().Debug("waiting for a VM to stop", zap.String("name", vmName), zap.String("namespace", vmNamespace))
			return false, nil
		}

		switch t := err.(type) {
		case *errors.StatusError:
			if t.Status().Reason == v1.StatusReasonNotFound {
				return true, nil
			}
			return false, err
		default:
			return false, err
		}

	})
}

func (e *Executor) EnsureVMDeleted() error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()

	return wait.PollImmediateInfinite(constants.PollVMtoDeleteInterval, func() (bool, error) {
		_, err := e.kubevirtClient.VirtualMachine(vmNamespace).Get(context.TODO(), vmName, &v1.GetOptions{})

		if err == nil {
			if err := e.ensureVMDelete(); err != nil {
				return false, err
			}
			log.Logger().Debug(" waiting for a VM to be deleted", zap.String("name", vmName), zap.String("namespace", vmNamespace))
			return false, nil
		}

		switch t := err.(type) {
		case *errors.StatusError:
			if t.Status().Reason == v1.StatusReasonNotFound {
				return true, nil
			}
			return false, err
		default:
			return false, err
		}

	})
}

func (e *Executor) SetupConnection(timeout time.Duration) error {
	if e.executor == nil {
		return fmt.Errorf("executor is missing or was not initialized")
	}

	if err := e.executor.Init(e.ipAddress); err != nil {
		return err
	}

	conditionFn := func() (done bool, err error) {
		return e.executor.TestConnection(), nil
	}

	var err error
	if timeout <= 0 {
		err = wait.PollImmediateInfinite(constants.PollValidConnectionInterval, conditionFn)
	} else {
		err = wait.PollImmediate(constants.PollValidConnectionInterval, timeout, conditionFn)
	}
	time.Sleep(constants.SetupConnectionDelay)
	return err
}

func (e *Executor) RemoteExecute(timeout time.Duration) error {
	if e.executor == nil {
		return fmt.Errorf("executor is missing or was not initialized")
	}
	return e.executor.RemoteExecute(timeout)
}

func (e *Executor) ensureVMStarted() error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()
	if !e.attemptedStart {
		e.attemptedStart = true
		log.Logger().Debug("starting a vm", zap.String("name", vmName), zap.String("namespace", vmNamespace))
		if err := e.kubevirtClient.VirtualMachine(vmNamespace).Start(context.TODO(), vmName, &kubevirtv1.StartOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) ensureVMStop() error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()
	if !e.attemptedStop {
		e.attemptedStop = true

		log.Logger().Debug("stopping a vm", zap.String("name", vmName), zap.String("namespace", vmNamespace))
		if err := e.kubevirtClient.VirtualMachine(vmNamespace).Stop(context.TODO(), vmName, &kubevirtv1.StopOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) ensureVMDelete() error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()
	if !e.attemptedDelete {
		e.attemptedDelete = true

		log.Logger().Debug("deleting a vm", zap.String("name", vmName), zap.String("namespace", vmNamespace))
		if err := e.kubevirtClient.VirtualMachine(vmNamespace).Delete(context.TODO(), vmName, &v1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}
