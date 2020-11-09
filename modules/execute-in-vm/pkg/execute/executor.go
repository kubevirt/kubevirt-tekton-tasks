package execute

import (
	"fmt"
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
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	"kubevirt.io/client-go/kubecli"
)

type Executor struct {
	clioptions     *parse.CLIOptions
	kubevirtClient kubecli.KubevirtClient
	executor       RemoteExecutor

	attemptedStart bool
	ipAddress      string
}

func NewExecutor(clioptions *parse.CLIOptions, execAttributes execattributes.ExecAttributes) (*Executor, error) {
	var executor RemoteExecutor

	switch execAttributes.GetType() {
	case constants.SSHSecretType:
		executor = newSSHExecutor(clioptions, execAttributes)
	default:
		return nil, fmt.Errorf("invalid secret/execution type %v", execAttributes.GetType())
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

func (e *Executor) EnsureVMRunning() error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()
	logFields := []zap.Field{zap.String("name", vmName), zap.String("namespace", vmNamespace)}

	//
	return wait.PollImmediateInfinite(constants.PollVMIInterval, func() (done bool, err error) {
		vmInstance, err := e.kubevirtClient.VirtualMachineInstance(vmNamespace).Get(vmName, &v1.GetOptions{})

		if err != nil {
			log.GetLogger().Debug("could not obtain a vm instance", logFields[0], logFields[1], zap.Error(err))
			switch t := err.(type) {
			case *errors.StatusError:
				if t.Status().Reason == v1.StatusReasonNotFound {
					if err := e.ensureVMStarted(); err != nil {
						return false, err
					}
					log.GetLogger().Debug(" waiting for a VMI to start", logFields...)
					return false, nil
				}
				return false, err
			default:
				return false, err
			}
		}

		switch vmInstance.Status.Phase {
		case kubevirtv1.Failed:
			log.GetLogger().Debug("vm instance failed", logFields[0], logFields[1], zap.Reflect("status", vmInstance.Status))
			// maybe the vm just stopped so let's try to start it again
			if err := e.ensureVMStarted(); err != nil {
				return false, err
			}
			log.GetLogger().Debug("waiting for a VMI to recover", logFields...)
			return false, nil
		case kubevirtv1.Running:
			ipAddress, ipError := vmi.GetPodIPAddress(vmInstance)

			if ipAddress == "" || ipError != nil {
				log.GetLogger().Debug("ip address not found", logFields[0], logFields[1], zap.Reflect("status", vmInstance.Status))

				if ipError != nil {
					return false, zerrors.NewMissingRequiredError(ipError.Error())
				}
				// wait for ipAddress
				return false, nil
			}
			log.GetLogger().Debug("ip address found", zap.String("ipAddress", ipAddress))
			e.ipAddress = ipAddress

			return true, nil

		default:
			log.GetLogger().Debug("waiting for a VMI to start", logFields...)
			return false, nil
		}
	})
}

func (e *Executor) SetupConnection() error {
	if err := e.executor.Init(e.ipAddress); err != nil {
		return err
	}

	return wait.PollImmediateInfinite(constants.PollValidConnectionInterval, func() (done bool, err error) {
		return e.executor.TestConnection(), nil
	})
}

func (e *Executor) RemoteExecute() error {
	return e.executor.RemoteExecute()
}

func (e *Executor) ensureVMStarted() error {
	vmName := e.clioptions.VirtualMachineName
	vmNamespace := e.clioptions.GetVirtualMachineNamespace()
	if !e.attemptedStart {
		e.attemptedStart = true
		log.GetLogger().Debug("starting a vm", zap.String("name", vmName), zap.String("namespace", vmNamespace))
		if err := e.kubevirtClient.VirtualMachine(vmNamespace).Start(vmName); err != nil {
			return err
		}
	}
	return nil
}
