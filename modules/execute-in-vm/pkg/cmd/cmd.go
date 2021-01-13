package cmd

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"os/exec"
	"time"
)

func RunCmdWithTimeout(timeout time.Duration, cmd *exec.Cmd) error {
	if timeout <= 0 {
		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return exit.Exit{
					Code: exitErr.ExitCode(),
					Soft: true,
				}
			} else {
				return err
			}
		}
	} else {
		if err := cmd.Start(); err != nil {
			return err
		}

		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		timeout := time.After(timeout)

		select {
		case <-timeout:
			cmd.Process.Kill()
			return exit.Exit{
				Code: constants.CommandTimeout,
				Msg:  "command timed out",
				Soft: true,
			}
		case err := <-done:
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					return exit.Exit{
						Code: exitErr.ExitCode(),
						Soft: true,
					}
				} else {
					return err
				}
			}
		}
	}

	return exit.Exit{
		Code: 0,
		Soft: true,
	}
}
