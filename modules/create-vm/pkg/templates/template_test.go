package templates_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
)

var _ = Describe("Template", func() {

	It("GetFlagLabelByPrefix", func() {
		key, value := templates.GetFlagLabelByPrefix(template.NewFedoraServerTinyTemplate().Build(), "workload.template.kubevirt.io")
		Expect(key).To(Equal("workload.template.kubevirt.io/server"))
		Expect(value).To(Equal("true"))
	})

	It("GetOs", func() {
		osID, osName := templates.GetOs(template.NewFedoraServerTinyTemplate().Build())
		Expect(osID).To(Equal("fedora29"))
		Expect(osName).To(Equal("Fedora 27 or higher"))
	})
})
