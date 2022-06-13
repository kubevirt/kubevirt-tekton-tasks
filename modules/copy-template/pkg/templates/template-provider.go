package templates

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	templatev1 "github.com/openshift/api/template/v1"
	v1 "github.com/openshift/api/template/v1"
	tempclient "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	templateclientset "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
)

type templateProvider struct {
	client tempclient.TemplateV1Interface
}

type TemplateProvider interface {
	Get(string, string) (*templatev1.Template, error)
	Create(*templatev1.Template) (*templatev1.Template, error)
	Update(*templatev1.Template) (*templatev1.Template, error)
}

func NewTemplateProvider(client tempclient.TemplateV1Interface) TemplateProvider {
	return &templateProvider{
		client: client,
	}
}

func (t *templateProvider) Get(namespace string, name string) (*templatev1.Template, error) {
	return t.client.Templates(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (t *templateProvider) Create(template *templatev1.Template) (*templatev1.Template, error) {
	return t.client.Templates(template.Namespace).Create(context.TODO(), template, metav1.CreateOptions{})
}

func (t *templateProvider) Update(template *templatev1.Template) (*templatev1.Template, error) {
	return t.client.Templates(template.Namespace).Update(context.TODO(), template, metav1.UpdateOptions{})
}

type TemplateCreator struct {
	cliOptions       *parse.CLIOptions
	templateProvider TemplateProvider
}

func NewTemplateCreator(cliOptions *parse.CLIOptions) (*TemplateCreator, error) {
	log.Logger().Debug("initialized clients and providers")

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return &TemplateCreator{
		cliOptions:       cliOptions,
		templateProvider: NewTemplateProvider(templateclientset.NewForConfigOrDie(config)),
	}, nil
}

func (t *TemplateCreator) CopyTemplate() (*v1.Template, error) {
	log.Logger().Debug("retrieving template", zap.String("name", t.cliOptions.GetSourceTemplateName()), zap.String("namespace", t.cliOptions.GetSourceTemplateNamespace()))
	template, err := t.templateProvider.Get(t.cliOptions.GetSourceTemplateNamespace(), t.cliOptions.GetSourceTemplateName())
	if err != nil {
		return nil, err
	}

	log.Logger().Debug("Original template metadata", zap.Any("ObjectMeta", template.ObjectMeta))

	if isCommonTemplate(template) {
		removeCommonTemplateInformationFromTemplate(template.Labels)
		removeCommonTemplateInformationFromTemplate(template.Annotations)

		unstructuredVM, err := getUnstructuredVM(template)
		if err != nil {
			return nil, err
		}

		if !(unstructuredVM.GetAPIVersion() == "kubevirt.io/v1" && unstructuredVM.GetKind() == "VirtualMachine") {
			return nil, fmt.Errorf("template %s contains unexpected object: %s, %s", template.Name, unstructuredVM.GetAPIVersion(), unstructuredVM.GetKind())
		}

		t.UpdateVMMetadata(unstructuredVM)

		t.EncodeVMToTemplate(template, unstructuredVM)
	}

	updatedTemplate := t.UpdateTemplateMetadata(template)

	log.Logger().Debug("Updated template metadata", zap.Any("ObjectMeta", updatedTemplate.ObjectMeta))
	existingTemplate, err := t.templateProvider.Get(t.cliOptions.GetTargetTemplateNamespace(), t.cliOptions.GetTargetTemplateName())

	if t.cliOptions.GetAllowReplaceValue() && existingTemplate != nil && err == nil {
		updatedTemplate.ResourceVersion = existingTemplate.ResourceVersion
		return t.templateProvider.Update(updatedTemplate)
	}

	return t.templateProvider.Create(updatedTemplate)
}
func removeCommonTemplateMetadataFromUnstructuredVM(unstructuredVM *unstructured.Unstructured, path []string, additionalMetadata map[string]string) error {
	obj, foundObj, err := unstructured.NestedStringMap(unstructuredVM.UnstructuredContent(), path...)
	if err != nil {
		return err
	}
	if foundObj {
		removeCommonTemplateInformationFromObj(obj)
		for key, value := range additionalMetadata {
			obj[key] = value
		}

		err := unstructured.SetNestedStringMap(unstructuredVM.UnstructuredContent(), obj, path...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TemplateCreator) UpdateVMMetadata(unstructuredVM *unstructured.Unstructured) error {
	labelPath := []string{"metadata", "labels"}
	err := removeCommonTemplateMetadataFromUnstructuredVM(unstructuredVM, labelPath, map[string]string{VMTemplateNameLabel: t.cliOptions.TargetTemplateName})
	if err != nil {
		return err
	}

	templateLabelPath := []string{"spec", "template", "metadata", "labels"}

	err = removeCommonTemplateMetadataFromUnstructuredVM(unstructuredVM, templateLabelPath, nil)
	if err != nil {
		return err
	}

	templateAnnotationsPath := []string{"spec", "template", "metadata", "annotations"}
	err = removeCommonTemplateMetadataFromUnstructuredVM(unstructuredVM, templateAnnotationsPath, nil)
	if err != nil {
		return err
	}

	annotationsPath := []string{"metadata", "annotations"}
	err = removeCommonTemplateMetadataFromUnstructuredVM(unstructuredVM, annotationsPath, nil)
	if err != nil {
		return err
	}

	return nil
}

func (t *TemplateCreator) EncodeVMToTemplate(template *templatev1.Template, unstructuredVM *unstructured.Unstructured) (*v1.Template, error) {
	raw, err := unstructuredVM.MarshalJSON()
	if err != nil {
		return nil, err
	}

	template.Objects[0].Raw = raw
	return template, nil
}

func (t *TemplateCreator) UpdateTemplateMetadata(template *v1.Template) *v1.Template {
	//set "template.kubevirt.io/type" label to VM so it is visible in UI
	template.Labels[TemplateTypeLabel] = VMTypeLabelValue

	newObjectMeta := metav1.ObjectMeta{
		Namespace:   t.cliOptions.GetTargetTemplateNamespace(),
		Labels:      template.Labels,
		Annotations: template.Annotations,
	}

	if t.cliOptions.GetTargetTemplateName() == "" {
		newObjectMeta.GenerateName = t.cliOptions.GetSourceTemplateName()
	} else {
		newObjectMeta.Name = t.cliOptions.GetTargetTemplateName()
	}

	template.ObjectMeta = newObjectMeta
	return template
}

func isCommonTemplate(template *v1.Template) bool {
	if val, ok := template.Labels[TemplateTypeLabel]; ok && val == templateTypeBaseValue {
		return true
	}
	return false
}

func removeCommonTemplateInformationFromTemplate(obj map[string]string) {
	for record, _ := range obj {
		if strings.HasPrefix(record, TemplateOsLabelPrefix) {
			delete(obj, record)
		}

		if strings.HasPrefix(record, TemplateFlavorLabelPrefix) {
			delete(obj, record)
		}

		if strings.HasPrefix(record, TemplateWorkloadLabelPrefix) {
			delete(obj, record)
		}
	}
	delete(obj, TemplateVersionLabel)
	delete(obj, TemplateDeprecatedAnnotation)
	delete(obj, KubevirtDefaultOSVariant)

	delete(obj, OpenshiftDocURL)
	delete(obj, OpenshiftProviderDisplayName)
	delete(obj, OpenshiftSupportURL)

	delete(obj, TemplateKubevirtProvider)
	delete(obj, TemplateKubevirtProviderSupportLevel)
	delete(obj, TemplateKubevirtProviderURL)

	delete(obj, OperatorSDKPrimaryResource)
	delete(obj, OperatorSDKPrimaryResourceType)

	delete(obj, AppKubernetesComponent)
	delete(obj, AppKubernetesName)
	delete(obj, AppKubernetesPartOf)
	delete(obj, AppKubernetesVersion)
	delete(obj, AppKubernetesManagedBy)
}

func removeCommonTemplateInformationFromObj(obj map[string]string) {
	delete(obj, VMFlavorAnnotation)
	delete(obj, VMOSAnnotation)
	delete(obj, VMWorkloadAnnotation)
	delete(obj, VMDomainLabel)
	delete(obj, VMSizeLabel)
	delete(obj, VMTemplateRevisionLabel)
	delete(obj, VMTemplateVersionLabel)
}

func getUnstructuredVM(template *templatev1.Template) (*unstructured.Unstructured, error) {
	unstructuredVM := &unstructured.Unstructured{}
	err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(template.Objects[0].Raw), 1024).Decode(unstructuredVM)
	if err != nil {
		return nil, err
	}
	return unstructuredVM, nil
}
