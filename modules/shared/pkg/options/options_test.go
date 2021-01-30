package options_test

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	input = " -a   1    -b \"hello world\" -u=5 -o 8 --long-d='rara=\"85\" ' --long-o=9 positional   arguments   "

	badInput = "-d=\" -c \"8\" "
)

var result = []string{"-a", "1", "-b", "hello world", "-u=5", "-o", "8", "--long-d=rara=\"85\" ", "--long-o=9", "positional", "arguments"}

var _ = Describe("Options", func() {
	It("operations works correctly", func() {
		opts, err := options.NewCommandOptions(input)
		fmt.Printf("%v", opts.GetAll())
		Expect(err).Should(Succeed())
		Expect(opts.Includes("-a")).To(BeTrue())
		Expect(opts.Includes("-o")).To(BeTrue())
		Expect(opts.Includes("-d")).To(BeFalse())
		Expect(opts.Includes("--long-d")).To(BeTrue())
		Expect(opts.GetOptionValue("-a")).To(Equal("1"))
		Expect(opts.GetOptionValue("-b")).To(Equal("hello world"))
		Expect(opts.GetOptionValue("-u")).To(Equal("5"))
		Expect(opts.GetOptionValue("--long-d")).To(Equal("rara=\"85\" "))
		Expect(opts.GetOptionValue("--not-included")).To(Equal(""))
		Expect(opts.GetOptionValue("-n")).To(Equal(""))
		Expect(opts.GetAll()).Should(Equal(result))

		opts.AddOpt("-d", "test")
		Expect(opts.Includes("-d")).To(BeTrue())
		Expect(opts.Includes("test")).To(BeFalse())
		latest := opts.GetAll()
		latestLen := len(latest)
		Expect(latest[latestLen-2]).Should(Equal("-d"))
		Expect(latest[latestLen-1]).Should(Equal("test"))
		Expect(opts.GetOptionValue("-d")).To(Equal("test"))

		opts.AddFlag("--verbose")
		Expect(opts.Includes("--verbose")).To(BeTrue())
		Expect(opts.GetOptionValue("--verbose")).To(Equal(""))
	})

	It("bad input", func() {
		opts, err := options.NewCommandOptions(badInput)
		Expect(opts).Should(BeNil())
		Expect(err).Should(HaveOccurred())
	})

	It("to string", func() {
		opts, err := options.NewCommandOptions("-o aaa -o \"bbbb cccc\" dd")
		Expect(err).Should(Succeed())
		Expect(opts.ToString()).Should(Equal("[-o, aaa, -o, bbbb cccc, dd]"))

		var nilOpts *options.CommandOptions
		Expect(nilOpts.ToString()).To(Equal("nil"))
	})

})
