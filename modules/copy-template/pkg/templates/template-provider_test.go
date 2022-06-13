package templates_test

import (
	"bytes"
	"encoding/json"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/templates"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift/api/template/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Template provider", func() {
	Context("Common templates information removed", func() {
		It("should remove Common template information", func() {
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
			updatedTemplate := tProvider.UpdateTemplateMetadata(t)
			Expect(t.GetLabels()).To(HaveLen(1))
			Expect(t.GetAnnotations()).To(HaveLen(1))

			for key, val := range t.GetLabels() {
				if key == "someOtherLabel" {
					Expect(updatedTemplate.Labels[key]).To(Equal(val))
				} else {
					Expect(updatedTemplate.Labels).ToNot(HaveKey(key))
				}
			}
			for key, val := range t.GetAnnotations() {
				if key == "someOtherLabel" {
					Expect(updatedTemplate.Labels[key]).To(Equal(val))
				} else {
					Expect(updatedTemplate.Labels).ToNot(HaveKey(key))
				}
			}

		})

		It("should remove Common template information from VM", func() {
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
			vmRaw, err := json.Marshal(vm)
			Expect(err).ToNot(HaveOccurred())

			unstructuredVM := &unstructured.Unstructured{}
			err = yaml.NewYAMLOrJSONDecoder(bytes.NewReader(vmRaw), 1024).Decode(unstructuredVM)
			Expect(err).ToNot(HaveOccurred())

			err = tProvider.UpdateVMMetadata(unstructuredVM)
			Expect(err).ToNot(HaveOccurred())

			labels, foundLabels, err := unstructured.NestedStringMap(unstructuredVM.UnstructuredContent(), []string{"metadata", "labels"}...)
			Expect(err).ToNot(HaveOccurred())
			Expect(foundLabels).To(BeTrue())
			Expect(labels).To(HaveLen(1))

			annotations, foundAnnotations, err := unstructured.NestedStringMap(unstructuredVM.UnstructuredContent(), []string{"metadata", "annotations"}...)
			Expect(err).ToNot(HaveOccurred())
			Expect(foundAnnotations).To(BeTrue())
			Expect(annotations).To(HaveLen(1))

			for key, val := range labels {
				if key == "someOtherLabel" {
					Expect(vm.Labels[key]).To(Equal(val))
				} else {
					Expect(vm.Labels).ToNot(HaveKey(key))
				}
			}
			for key, val := range annotations {
				if key == "someOtherLabel" {
					Expect(vm.Labels[key]).To(Equal(val))
				} else {
					Expect(vm.Labels).ToNot(HaveKey(key))
				}
			}

		})
	})
})
