package intutilstest

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
)

func SetupTestSuite() {
	utilstest.SetEnv(zconstants.OutOfClusterENV, "true")
}

func TearDownSuite() {
	utilstest.UnSetEnv(zconstants.OutOfClusterENV)
}
