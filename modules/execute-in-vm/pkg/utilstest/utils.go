package utilstest

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/logger"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"go.uber.org/zap"
)

func SetupTestSuite() {
	utilstest.SetEnv(zconstants.OutOfClusterENV, "true")
	logger.InitLogger(zap.InfoLevel)
}

func TearDownSuite() {
	utilstest.UnSetEnv(zconstants.OutOfClusterENV)
}
