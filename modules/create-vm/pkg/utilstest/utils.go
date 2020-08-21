package utilstest

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/log"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"go.uber.org/zap"
)

func SetupTestSuite() {
	utilstest.SetEnv(zconstants.OutOfClusterENV, "true")
	log.InitLogger(zap.InfoLevel)
}

func TearDownSuite() {
	utilstest.UnSetEnv(zconstants.OutOfClusterENV)
}
