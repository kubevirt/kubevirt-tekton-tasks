package requirements

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

func GetLabelRequirement(condition string) (labels.Requirements, error) {
	if strings.TrimSpace(condition) == "" {
		return nil, nil
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
