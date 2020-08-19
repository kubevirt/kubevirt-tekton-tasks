package validations

import (
	"strings"
)

const (
	diskBusJSONPath = "jsonpath::.spec.domain.devices.disks[*].disk.bus"
)

const (
	diskBusVirtIO = "virtio"
	diskBusSATA   = "sata"
	diskBusSCSI   = "scsi"
)

var defaultDiskBus = diskBusVirtIO

type CommonTemplateValidation struct {
	Name        string   // Identifier of the rule. Must be unique among all the rules attached to a template
	Rule        string   // Validation rule name. One of integer|string|regex|enum
	Path        string   // jsonpath of the field whose value is going to be evaluated.
	Message     string   // User-friendly string message describing the failure, should the rule not be satisfied
	Min         int      // For 'integer' rule
	Max         int      // For 'integer' rule
	MinLength   int      // For 'string' rule
	MaxLength   int      // For 'string' rule
	Regex       string   // For 'regex' rule
	Values      []string // For 'enum' rule
	JustWarning bool
}

type TemplateValidations struct {
	validations []CommonTemplateValidation
}

func NewTemplateValidations(validations []CommonTemplateValidation) *TemplateValidations {
	return &TemplateValidations{validations}
}

func (t *TemplateValidations) IsEmpty() bool {
	return len(t.validations) == 0
}

func (t *TemplateValidations) GetDefaultDiskBus() string {

	allowedBuses := t.getAllowedBuses(diskBusJSONPath, false)

	if len(allowedBuses) == 0 {
		return defaultDiskBus
	}

	recommendedBuses := t.getRecommendedBuses(diskBusJSONPath)

	if len(recommendedBuses) > 0 {
		if recommendedBuses[defaultDiskBus] {
			return defaultDiskBus
		}
		for bus := range recommendedBuses {
			return bus
		}
	}

	if allowedBuses[defaultDiskBus] {
		return defaultDiskBus
	}
	for bus := range allowedBuses {
		return bus
	}

	return defaultDiskBus
}

func (t *TemplateValidations) getAllowedBuses(jsonPath string, justWarning bool) map[string]bool {
	allowedBuses := t.getAllowedEnumValues(jsonPath, justWarning)
	if len(allowedBuses) == 0 {
		// default to all
		return map[string]bool{
			diskBusVirtIO: true,
			diskBusSATA:   true,
			diskBusSCSI:   true,
		}
	}

	result := make(map[string]bool)
	for _, diskBus := range allowedBuses {
		result[diskBus] = true
	}
	return result
}

func (t *TemplateValidations) getRecommendedBuses(jsonPath string) map[string]bool {
	allowedBuses := t.getAllowedBuses(jsonPath, false)
	recommendedBuses := t.getAllowedBuses(jsonPath, true)

	for key := range recommendedBuses {
		if !allowedBuses[key] {
			delete(recommendedBuses, key)
		}
	}

	if len(recommendedBuses) == 0 {
		return allowedBuses
	}

	return recommendedBuses
}

func (t *TemplateValidations) getAllowedEnumValues(jsonPath string, justWarning bool) []string {
	var relevantValidations []string

	for _, validation := range t.validations {
		if strings.HasPrefix(validation.Path, jsonPath) && validation.JustWarning == justWarning {
			relevantValidations = append(relevantValidations, validation.Values...)
		}
	}
	return relevantValidations
}
