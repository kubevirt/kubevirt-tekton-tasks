package options_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	. "github.com/onsi/ginkgo/v2"
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
		Expect(err).Should(Succeed())
		Expect(opts.IncludesOption("-a")).To(BeTrue())
		Expect(opts.IncludesOption("-o")).To(BeTrue())
		Expect(opts.IncludesOption("-d")).To(BeFalse())
		Expect(opts.IncludesOption("--long-d")).To(BeTrue())
		Expect(opts.IncludesString("hello world")).To(BeTrue())
		Expect(opts.IncludesString("-o=9 positional arguments")).To(BeTrue())
		Expect(opts.GetOptionValue("-a")).To(Equal("1"))
		Expect(opts.GetOptionValue("-b")).To(Equal("hello world"))
		Expect(opts.GetOptionValue("-u")).To(Equal("5"))
		Expect(opts.GetOptionValue("--long-d")).To(Equal("rara=\"85\" "))
		Expect(opts.GetOptionValue("--not-included")).To(Equal(""))
		Expect(opts.GetOptionValue("-n")).To(Equal(""))
		Expect(opts.GetAll()).Should(Equal(result))

		opts.AddOption("-d", "test")
		Expect(opts.IncludesOption("-d")).To(BeTrue())
		Expect(opts.IncludesOption("test")).To(BeFalse())
		Expect(opts.IncludesString("test")).To(BeTrue())
		latest := opts.GetAll()
		latestLen := len(latest)
		Expect(latest[latestLen-2]).Should(Equal("-d"))
		Expect(latest[latestLen-1]).Should(Equal("test"))
		Expect(opts.GetOptionValue("-d")).To(Equal("test"))

		opts.AddFlag("--verbose")
		Expect(opts.IncludesOption("--verbose")).To(BeTrue())
		Expect(opts.IncludesString("--verbose")).To(BeTrue())
		Expect(opts.GetOptionValue("--verbose")).To(Equal(""))

		opts.AddValue("false")
		Expect(opts.GetOptionValue("--verbose")).To(Equal("false"))

		opts.AddOptions("1", "2", "3")
		Expect(opts.IncludesString("1 2 3")).To(BeTrue())
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

	It("second constructor works", func() {
		opts := options.NewCommandOptionsFromArray([]string{"-a", "b", "-c=d"})
		Expect(opts.ToString()).Should(Equal("[-a, b, -c=d]"))
	})

	It("test short options", func() {
		opts, err := options.NewCommandOptions("-a=b -c --dery-long arg -d e -fthis --long-option -u\"256B\" -osshStrictHostKeyCheckingOption=true --pp p50 -p40")
		Expect(err).Should(Succeed())

		for optKey, optVal := range map[string]string{
			"-a":            "b",
			"-c":            "",
			"--dery-long":   "arg",
			"-d":            "e",
			"-f":            "this",
			"--long-option": "",
			"-u":            "256B",
			"-o":            "sshStrictHostKeyCheckingOption=true",
			"-p":            "40",
		} {
			Expect(opts.IncludesOption(optKey)).To(BeTrue())
			Expect(opts.GetOptionValue(optKey)).To(Equal(optVal))
			Expect(opts.IncludesString(optKey)).To(BeTrue())
			Expect(opts.IncludesString(optVal)).To(BeTrue())
		}
	})
})
