package utilstest

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	. "github.com/onsi/gomega"
	"os"
)

// logic should equal to utils.HandleExit func
func HandleTestExit(shouldPanic bool, shouldExitWithCode int, shouldExitWithMessage string) {
	if e := recover(); e != nil {
		if ex, ok := e.(exit.Exit); ok == true && ex.Soft {
			errMsg := ex.Error()
			if len(errMsg) > 0 {
				if errMsg[len(errMsg)-1] != '\n' {
					errMsg += "\n"
				}
			}
			Expect(errMsg).To(Equal(shouldExitWithMessage))
			Expect(ex.Code).To(Equal(shouldExitWithCode))
			Expect(shouldPanic).To(BeFalse(), "should not panic")
			return
		}
		Expect(shouldPanic).To(BeTrue(), "should panic")
	}
}

func SetEnv(key, value string) {
	Expect(os.Setenv(key, value)).Should(Succeed())
}

func UnSetEnv(key string) {
	Expect(os.Unsetenv(key)).Should(Succeed())
}
