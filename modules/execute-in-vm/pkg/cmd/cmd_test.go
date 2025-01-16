package cmd_test

import (
	"os/exec"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cmd", func() {

	Describe("RunCmdWithTimeout", func() {
		DescribeTable("finish without reaching timeout", func(command *exec.Cmd, timeout time.Duration, shouldRun bool) {
			start := time.Now()
			err := cmd.RunCmdWithTimeout(timeout, command)
			runDuration := time.Since(start)
			if shouldRun {
				Expect(err).ShouldNot(HaveOccurred())
			} else {
				Expect(err).Should(HaveOccurred())
			}

			if timeout > 0 {
				Expect(runDuration).Should(BeNumerically("<=", timeout))
			}
		},
			Entry("unlimited timeout", exec.Command("echo", "-n"), 0*time.Second, true),
			Entry("unlimited timeout fail", exec.Command("false"), 0*time.Second, false),
			Entry("unlimited timeout invalid cmd", exec.Command("invalidcmd"), 0*time.Second, false),
			//
			Entry("3 sec timeout", exec.Command("echo", "-n"), 3*time.Second, true),
			Entry("3 sec timeout fail", exec.Command("false"), 3*time.Second, false),
			Entry("3 sec timeout invalid cmd", exec.Command("invalidcmd"), 3*time.Second, false),
		)

		It("times out", func() {
			timeout := 100 * time.Millisecond
			start := time.Now()
			err := cmd.RunCmdWithTimeout(timeout, exec.Command("sleep", "1"))
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("timed out"))
			Expect(time.Since(start)).Should(BeNumerically(">=", timeout))
		})

	})
})
