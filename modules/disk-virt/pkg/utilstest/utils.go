package utilstest

import (
	"go.uber.org/zap"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/log"
)

func SetupTestSuite() {
	log.InitLogger(zap.InfoLevel)
}
