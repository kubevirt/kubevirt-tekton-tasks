package templates

import (
	"encoding/json"
	templatev1 "github.com/openshift/api/template/v1"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	lab "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants/labels"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/templates/validations"
	"k8s.io/apimachinery/pkg/runtime"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	"sort"
	"strings"
)

const (
	validationsAnnotation = "validations"
	osLabelPrefix         = lab.TemplateOsLabel + "/"
)

// Get label whose key starts with labelPrefix and has value true
// returns key, val if such label is found
func GetFlagLabelByPrefix(template *templatev1.Template, labelPrefix string) (string, string) {
	if labels := template.GetLabels(); labels != nil {
		for key, val := range labels {
			if strings.HasPrefix(key, labelPrefix) && val == "true" {
				return key, val
			}
		}
	}
	return "", ""
}

func DecodeVM(template *templatev1.Template) (*kubevirtv1.VirtualMachine, error) {
	var vm = &kubevirtv1.VirtualMachine{}

	for _, obj := range template.Objects {
		decoder := kubevirtv1.Codecs.UniversalDecoder(kubevirtv1.GroupVersion)
		decoded, err := runtime.Decode(decoder, obj.Raw)
		if err != nil {
			return nil, err
		}
		done, ok := decoded.(*kubevirtv1.VirtualMachine)
		if ok {
			vm = done
			break
		}
	}
	return vm, nil
}

func GetTemplateValidations(template *templatev1.Template) (*validations.TemplateValidations, error) {
	marshalledValidations := template.Annotations[validationsAnnotation]
	var commonTemplateValidations []validations.CommonTemplateValidation

	// empty validations have defaults
	if marshalledValidations != "" {
		if err := json.Unmarshal([]byte(marshalledValidations), &commonTemplateValidations); err != nil {
			return nil, err
		}
	}
	return validations.NewTemplateValidations(commonTemplateValidations), nil
}

// returns osID, osName
func GetOs(template *templatev1.Template) (string, string) {

	var osIds textIDs

	for key, val := range template.Labels {
		if strings.HasPrefix(key, osLabelPrefix) && val == constants.True {
			osId := key[len(osLabelPrefix):]
			osIds = append(osIds, osId)
		}
	}

	sort.Sort(osIds)

	if len(osIds) == 0 {
		return "", ""
	}

	osID := osIds[len(osIds)-1]

	osName := template.Annotations[lab.TemplateNameOsAnnotation+"/"+osID]

	return osID, osName
}
