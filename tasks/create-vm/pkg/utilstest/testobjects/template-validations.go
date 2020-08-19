package testobjects

import (
	"strings"

	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/templates/validations"
)

func NewTestCommonTemplateValidations(buses ...string) []validations.CommonTemplateValidation {
	if len(buses) == 0 {
		return []validations.CommonTemplateValidation{}
	}

	var values []string

	values = append(values, buses...)

	return []validations.CommonTemplateValidation{{
		Name:    "disk-bus",
		Rule:    "enum",
		Path:    "jsonpath::.spec.domain.devices.disks[*].disk.bus",
		Message: strings.Join(values, ", ") + " allowed only",
		Values:  values,
	}}
}
