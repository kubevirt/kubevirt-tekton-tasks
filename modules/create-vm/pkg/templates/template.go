package templates

import (
	"encoding/json"
	"sort"
	"strings"

	lab "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants/labels"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	templatev1 "github.com/openshift/api/template/v1"
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
		if strings.HasPrefix(key, osLabelPrefix) && val == zconstants.True {
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
