package vmi_test

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/vmi"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "kubevirt.io/client-go/api/v1"
)

var _ = Describe("VMI", func() {
	var vmi *v1.VirtualMachineInstance

	BeforeEach(func() {
		vmi = testobjects.NewTestVMI()
	})

	Describe("GetPodIPAddress", func() {
		It("returns empty by default", func() {
			ipAddress, err := GetPodIPAddress(vmi)
			Expect(err).To(BeNil())
			Expect(ipAddress).To(BeEmpty())
		})
		It("no pod network", func() {
			vmi.Spec.Networks = []v1.Network{}
			vmi.Spec.Domain.Devices.Interfaces = nil
			ipAddress, err := GetPodIPAddress(vmi)
			Expect(err).Should(HaveOccurred())
			Expect(ipAddress).Should(BeEmpty())
		})
		It("returns IP address", func() {
			ip := "135.21.75.16"
			vmi.Status = v1.VirtualMachineInstanceStatus{
				Interfaces: []v1.VirtualMachineInstanceNetworkInterface{
					{
						Name: vmi.Spec.Networks[0].Name,
						IP:   ip,
					},
				},
			}
			ipAddress, err := GetPodIPAddress(vmi)
			Expect(err).Should(Succeed())
			Expect(ipAddress).Should(Equal(ip))
		})
	})
})
