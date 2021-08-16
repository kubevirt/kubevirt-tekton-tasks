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
	Create(string, *templatev1.Template) (*templatev1.Template, error)
}

func NewTemplateProvider(client tempclient.TemplateV1Interface) TemplateProvider {
	return &templateProvider{
		client: client,
	}
}

func (t *templateProvider) Get(namespace string, name string) (*templatev1.Template, error) {
	return t.client.Templates(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (t *templateProvider) Create(namespace string, template *templatev1.Template) (*templatev1.Template, error) {
	return t.client.Templates(namespace).Create(context.TODO(), template, metav1.CreateOptions{})
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

	updatedTemplate, err := t.UpdateTemplate(template)
	if err != nil {
		return nil, err
	}

	return t.templateProvider.Create(t.cliOptions.GetTargetTemplateNamespace(), updatedTemplate)
}

func (t *TemplateCreator) UpdateTemplate(template *v1.Template) (*v1.Template, error) {
	newObjectMeta := metav1.ObjectMeta{
		Name:        t.cliOptions.GetTargetTemplateName(),
		Namespace:   t.cliOptions.GetTargetTemplateNamespace(),
		Labels:      template.Labels,
		Annotations: template.Annotations,
	}
	template.ObjectMeta = newObjectMeta
	return template, nil
}
