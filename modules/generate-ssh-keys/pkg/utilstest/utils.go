package utilstest

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/log"
	"go.uber.org/zap"
)

func SetupTestSuite() {
	log.InitLogger(zap.InfoLevel)
}
