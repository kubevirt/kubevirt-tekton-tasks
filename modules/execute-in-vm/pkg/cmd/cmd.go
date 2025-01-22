package cmd

import (
	"errors"
	"os/exec"
	"time"
)

func RunCmdWithTimeout(timeout time.Duration, cmd *exec.Cmd) error {
	if timeout <= 0 {
		if err := cmd.Run(); err != nil {
			return err
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
			return errors.New("command timed out")
		case err := <-done:
			if err != nil {
				return err
			}
		}
	}
	return nil
}
