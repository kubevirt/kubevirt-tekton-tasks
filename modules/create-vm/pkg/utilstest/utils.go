package utilstest

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"go.uber.org/zap"
)

func SetupTestSuite() {
	log.InitLogger(zap.InfoLevel)
}
