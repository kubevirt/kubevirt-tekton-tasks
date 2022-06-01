package templates_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/templates"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift/api/template/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Template provider", func() {
	Context("Common templates informations removed", func() {
		It("should remove Common template informations", func() {
			tProvider := &templates.TemplateCreator{}
			t := &v1.Template{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						templates.OpenshiftDocURL:                      "test",
						templates.OpenshiftProviderDisplayName:         "test",
						templates.OpenshiftSupportURL:                  "test",
						templates.KubevirtDefaultOSVariant:             "test",
						templates.TemplateKubevirtProvider:             "test",
						templates.TemplateKubevirtProviderSupportLevel: "test",
						templates.TemplateKubevirtProviderURL:          "test",
						templates.OperatorSDKPrimaryResource:           "test",
						templates.OperatorSDKPrimaryResourceType:       "test",
						templates.AppKubernetesComponent:               "test",
						templates.AppKubernetesManagedBy:               "test",
						templates.AppKubernetesName:                    "test",
						templates.AppKubernetesPartOf:                  "test",
						templates.AppKubernetesVersion:                 "test",
						templates.TemplateVersionLabel:                 "test",
						templates.TemplateTypeLabel:                    "test",
						templates.TemplateOsLabelPrefix:                "test",
						templates.TemplateFlavorLabelPrefix:            "test",
						templates.TemplateWorkloadLabelPrefix:          "test",
						templates.TemplateDeprecatedAnnotation:         "test",
						"someOtherLabel":                               "test",
					},
					Annotations: map[string]string{
						templates.OpenshiftDocURL:                      "test",
						templates.OpenshiftProviderDisplayName:         "test",
						templates.OpenshiftSupportURL:                  "test",
						templates.KubevirtDefaultOSVariant:             "test",
						templates.TemplateKubevirtProvider:             "test",
						templates.TemplateKubevirtProviderSupportLevel: "test",
						templates.TemplateKubevirtProviderURL:          "test",
						templates.OperatorSDKPrimaryResource:           "test",
						templates.OperatorSDKPrimaryResourceType:       "test",
						templates.AppKubernetesComponent:               "test",
						templates.AppKubernetesManagedBy:               "test",
						templates.AppKubernetesName:                    "test",
						templates.AppKubernetesPartOf:                  "test",
						templates.AppKubernetesVersion:                 "test",
						templates.TemplateVersionLabel:                 "test",
						templates.TemplateTypeLabel:                    "test",
						templates.TemplateOsLabelPrefix:                "test",
						templates.TemplateFlavorLabelPrefix:            "test",
						templates.TemplateWorkloadLabelPrefix:          "test",
						templates.TemplateDeprecatedAnnotation:         "test",
						"someOtherLabel":                               "test",
					},
				},
			}
			updatedTemplate := tProvider.UpdateTemplateMetaObject(t)
			Expect(len(t.GetLabels())).To(Equal(1))
			Expect(len(t.GetAnnotations())).To(Equal(1))

			for key, val := range t.GetLabels() {
				if key == "someOtherLabel" {
					Expect(updatedTemplate.Labels[key]).To(Equal(val))
				} else {
					Expect(updatedTemplate.Labels[key]).To(Equal(""))
				}
			}
			for key, val := range t.GetAnnotations() {
				if key == "someOtherLabel" {
					Expect(updatedTemplate.Labels[key]).To(Equal(val))
				} else {
					Expect(updatedTemplate.Labels[key]).To(Equal(""))
				}
			}

		})

		It("should remove Common template informations from VM", func() {
			tProvider := &templates.TemplateCreator{}

			vm := &kubevirtv1.VirtualMachine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						templates.VMFlavorAnnotation:   "test",
						templates.VMOSAnnotation:       "test",
						templates.VMWorkloadAnnotation: "test",
						templates.VMDomainLabel:        "test",
						templates.VMSizeLabel:          "test",
						"someOtherLabel":               "test",
					},
					Annotations: map[string]string{
						templates.VMFlavorAnnotation:   "test",
						templates.VMOSAnnotation:       "test",
						templates.VMWorkloadAnnotation: "test",
						templates.VMDomainLabel:        "test",
						templates.VMSizeLabel:          "test",
						"someOtherLabel":               "test",
					},
				},
			}

			tProvider.UpdateVMMetaObject(vm)
			Expect(len(vm.GetLabels())).To(Equal(1))
			Expect(len(vm.GetAnnotations())).To(Equal(1))

			for key, val := range vm.GetLabels() {
				if key == "someOtherLabel" {
					Expect(vm.Labels[key]).To(Equal(val))
				} else {
					Expect(vm.Labels[key]).To(Equal(""))
				}
			}
			for key, val := range vm.GetAnnotations() {
				if key == "someOtherLabel" {
					Expect(vm.Labels[key]).To(Equal(val))
				} else {
					Expect(vm.Labels[key]).To(Equal(""))
				}
			}

		})
	})
})
