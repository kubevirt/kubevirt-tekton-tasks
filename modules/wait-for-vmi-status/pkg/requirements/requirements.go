package requirements

import (
	"errors"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"k8s.io/apimachinery/pkg/labels"
)

func ValidateJSONPath(condition string) bool {
	// jsonpath format must be jsonpath='{.status.phase}' == Success"
	if strings.HasPrefix(condition, "jsonpath='{.") && strings.Contains(condition, "}'") {
		return true
	}

	return false
}

func ParseJSONPathParameter(condition string) (string, error) {
	if !ValidateJSONPath(condition) {
		return "", errors.New("valid jsonpath format is jsonpath='{.status.phase}' == Success")
	}

	condition = strings.Replace(condition, "jsonpath='{.", "", -1)
	return strings.Replace(condition, "}'", "", -1), nil
}

func GetLabelRequirement(condition string) (labels.Requirements, error) {
	var err error

	if strings.TrimSpace(condition) == "" {
		return nil, nil
	}

	if strings.HasPrefix(condition, "jsonpath=") {
		condition, err = ParseJSONPathParameter(condition)

		if err != nil {
			return nil, err
		}
	}

	selector, err := labels.Parse(condition)

	if err != nil {
		return nil, zerrors.NewMissingRequiredError("could not parse condition %v: %v", condition, err.Error())
	}

	requirements, selectable := selector.Requirements()

	if !selectable {
		return nil, nil
	}

	if _, err := ObjectToLabelsLookup(nil, requirements); err != nil {
		return nil, zerrors.NewMissingRequiredError("invalid condition: %v", err.Error())
	}
	return requirements, nil
}

// MatchesRequirements should be called after GetLabelRequirement to ensure consistency of jsonpath keys
func MatchesRequirements(obj interface{}, requirements labels.Requirements) bool {
	numMatched := 0

	// err checked in GetLabelRequirement
	objectMockingLabels, _ := ObjectToLabelsLookup(obj, requirements)

	for _, req := range requirements {
		matched := req.Matches(objectMockingLabels)
		if matched {
			numMatched++
		}
	}

	return numMatched == len(requirements)
}
