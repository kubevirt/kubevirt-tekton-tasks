package zutils

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"strings"
)

func IsTrue(value string) bool {
	return strings.ToLower(value) == zconstants.True
}
