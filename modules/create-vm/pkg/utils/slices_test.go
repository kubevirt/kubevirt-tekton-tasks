package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utils"
)

var _ = Describe("Slices", func() {

	Describe("Gets last", func() {
		It("does not", func() {
			Expect(utils.GetLast(nil)).To(BeEmpty())
			Expect(utils.GetLast([]string{})).To(BeEmpty())
		})
		It("gets", func() {
			Expect(utils.GetLast([]string{"a"})).To(Equal("a"))
			Expect(utils.GetLast([]string{"a", "b", "c"})).To(Equal("c"))
		})
	})

	Describe("Concat slices", func() {
		It("empty", func() {
			Expect(utils.ConcatStringSlices(nil, nil)).To(BeEmpty())
			Expect(utils.ConcatStringSlices([]string{}, nil)).To(BeEmpty())
		})
		It("concats", func() {
			concatenated := utils.ConcatStringSlices([]string{"a"}, []string{"b"})
			Expect(concatenated).To(HaveLen(2))
			Expect(concatenated).To(Equal([]string{"a", "b"}))
			Expect(utils.ConcatStringSlices(nil, []string{"a"})).To(HaveLen(1))
			Expect(utils.ConcatStringSlices([]string{"a"}, []string{"b", "c"})).To(HaveLen(3))
		})
	})

})
