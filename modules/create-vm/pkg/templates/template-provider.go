package templates

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	templatev1 "github.com/openshift/api/template/v1"
	tempclient "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	processingURI = "processedTemplates"
)

type templateProvider struct {
	client tempclient.TemplateV1Interface
}

type TemplateProvider interface {
	Get(namespace string, name string) (*templatev1.Template, error)
	Process(namespace string, template *templatev1.Template, paramValues map[string]string) (*templatev1.Template, error)
}

func NewTemplateProvider(client tempclient.TemplateV1Interface) TemplateProvider {
	return &templateProvider{
		client: client,
	}
}

func (t *templateProvider) Get(namespace string, name string) (*templatev1.Template, error) {
	return t.client.Templates(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (t *templateProvider) Process(namespace string, template *templatev1.Template, paramValues map[string]string) (*templatev1.Template, error) {
	temp := template.DeepCopy()
	temp.Namespace = ""
	params := temp.Parameters

	var paramsError zerrors.MultiError
	for i, param := range params {
		additionalValue := paramValues[param.Name]
		if additionalValue != "" {
			temp.Parameters[i].Value = additionalValue
		} else if param.Value == "" && param.Required {
			paramsError.Add(param.Name, zerrors.NewMissingRequiredError("required param %v is missing a value", param.Name))
		}
	}
	if !paramsError.IsEmpty() {
		return nil, paramsError.ShortPrint("required params are missing values:").AsOptional()
	}

	processedTemplate := &templatev1.Template{}
	err := t.client.RESTClient().Post().
		Namespace(namespace).
		Resource(processingURI).
		Body(temp).
		Do(context.TODO()).
		Into(processedTemplate)
	if err != nil {
		return nil, err
	}
	//setting namespace back for usage in VM labels
	processedTemplate.Namespace = template.Namespace
	return processedTemplate, nil
}
