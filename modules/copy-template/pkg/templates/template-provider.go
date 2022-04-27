package templates

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	templatev1 "github.com/openshift/api/template/v1"
	v1 "github.com/openshift/api/template/v1"
	tempclient "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	templateclientset "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	updatedTemplate := t.UpdateTemplateMetaObject(template)

	log.Logger().Debug("Updated template metadata", zap.Any("ObjectMeta", updatedTemplate.ObjectMeta))
	existingTemplate, err := t.templateProvider.Get(t.cliOptions.GetTargetTemplateNamespace(), t.cliOptions.GetTargetTemplateName())

	if t.cliOptions.GetAllowReplaceValue() && existingTemplate != nil && err == nil {
		updatedTemplate.ResourceVersion = existingTemplate.ResourceVersion
		return t.templateProvider.Update(updatedTemplate)
	}

	return t.templateProvider.Create(updatedTemplate)
}

func (t *TemplateCreator) UpdateTemplateMetaObject(template *v1.Template) *v1.Template {
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
