package cmd_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/cmd"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"os/exec"
	"time"
)

var _ = Describe("cmd", func() {

	Describe("RunCmdWithTimeout", func() {
		table.DescribeTable("finish without reaching timeout", func(command *exec.Cmd, timeout time.Duration, exitCode int, shouldRun bool) {
			start := time.Now()
			err := cmd.RunCmdWithTimeout(timeout, command)
			runDuration := time.Since(start)
			Expect(err).ShouldNot(BeNil())

			if shouldRun {
				if exitErr, ok := err.(exit.Exit); ok == true {
					Expect(exitErr.Code).To(Equal(exitCode))
				} else {
					Fail("did not return exit code:" + err.Error())
				}
			} else {
				_, ok := err.(exit.Exit)
				Expect(ok).To(BeFalse())
			}

			if timeout > 0 {
				Expect(runDuration).Should(BeNumerically("<=", timeout))
			}
		},
			table.Entry("unlimited timeout", exec.Command("echo", "-n"), 0*time.Second, 0, true),
			table.Entry("unlimited timeout fail", exec.Command("false"), 0*time.Second, 1, true),
			table.Entry("unlimited timeout invalid cmd", exec.Command("invalidcmd"), 0*time.Second, -1, false),
			//
			table.Entry("3 sec timeout", exec.Command("echo", "-n"), 3*time.Second, 0, true),
			table.Entry("3 sec timeout fail", exec.Command("false"), 3*time.Second, 1, true),
			table.Entry("3 sec timeout invalid cmd", exec.Command("invalidcmd"), 3*time.Second, -1, false),
		)

		It("times out", func() {
			timeout := 100 * time.Millisecond
			start := time.Now()
			err := cmd.RunCmdWithTimeout(timeout, exec.Command("sleep", "1"))
			Expect(err).ShouldNot(BeNil())
			exitCode := 0
			if exitErr, ok := err.(exit.Exit); ok == true {
				exitCode = exitErr.Code
			}
			Expect(err.Error()).To(ContainSubstring("timed out"))
			Expect(exitCode).To(Equal(constants.CommandTimeout))
			Expect(time.Since(start)).Should(BeNumerically(">=", timeout))
		})

	})
})
