package utils_test

import (
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VMI", func() {

	Describe("WithTimeout", func() {
		It("does not finish before timeout", func() {
			originalTimeout := 1 * time.Second
			runWithTimeout := utils.WithTimeout(originalTimeout)

			for i := 1; i < 3; i++ {
				called := false
				nextSleep := 10 * time.Millisecond
				time.Sleep(10 * time.Millisecond)

				runWithTimeout(func(timeout time.Duration, finished bool) {
					Expect(timeout).Should(BeNumerically("<", originalTimeout-(time.Duration(i)*nextSleep)))
					Expect(finished).Should(BeFalse())
					called = true
				})
				Expect(called).To(BeTrue())
			}
		})

		It("finishes after timeout", func() {
			originalTimeout := 10 * time.Millisecond
			runWithTimeout := utils.WithTimeout(originalTimeout)

			called := false
			nextSleep := 11 * time.Millisecond
			time.Sleep(nextSleep)

			runWithTimeout(func(timeout time.Duration, finished bool) {
				Expect(timeout).Should(BeNumerically("<", 0))
				Expect(finished).Should(BeTrue())
				called = true
			})
			Expect(called).To(BeTrue())
		})

		It("executes correctly without timeout", func() {
			runWithTimeout := utils.WithTimeout(0)

			for i := 0; i < 5; i++ {
				called := false
				time.Sleep(2 * time.Millisecond)
				runWithTimeout(func(timeout time.Duration, finished bool) {
					Expect(timeout).Should(BeZero())
					Expect(finished).Should(BeFalse())
					called = true
				})
				Expect(called).To(BeTrue())
			}
		})

	})
})
