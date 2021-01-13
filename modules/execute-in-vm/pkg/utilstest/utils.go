package utilstest

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"go.uber.org/zap"
)

func SetupTestSuite() {
	log.InitLogger(zap.InfoLevel)
}

func TearDownSuite() {
}
