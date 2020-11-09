package zutils_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Slices", func() {

	Describe("Gets last", func() {
		It("does not", func() {
			Expect(zutils.GetLast(nil)).To(BeEmpty())
			Expect(zutils.GetLast([]string{})).To(BeEmpty())
		})
		It("gets", func() {
			Expect(zutils.GetLast([]string{"a"})).To(Equal("a"))
			Expect(zutils.GetLast([]string{"a", "b", "c"})).To(Equal("c"))
		})
	})

	Describe("Concat slices", func() {
		It("empty", func() {
			Expect(zutils.ConcatStringSlices(nil, nil)).To(BeEmpty())
			Expect(zutils.ConcatStringSlices([]string{}, nil)).To(BeEmpty())
		})
		It("concats", func() {
			concatenated := zutils.ConcatStringSlices([]string{"a"}, []string{"b"})
			Expect(concatenated).To(HaveLen(2))
			Expect(concatenated).To(Equal([]string{"a", "b"}))
			Expect(zutils.ConcatStringSlices(nil, []string{"a"})).To(HaveLen(1))
			Expect(zutils.ConcatStringSlices([]string{"a"}, []string{"b", "c"})).To(HaveLen(3))
		})
	})

})
