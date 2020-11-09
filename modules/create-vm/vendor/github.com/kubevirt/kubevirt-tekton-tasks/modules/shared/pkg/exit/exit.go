package exit

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"os"
)

type Exit struct {
	Code int
	Msg  string
	Soft bool
}

func (e Exit) Error() string {
	return e.Msg
}

func (e Exit) IsSoft() bool {
	return e.Soft
}

func ExitFromError(code int, err error) {
	if err == nil {
		panic(Exit{
			Code: code,
			Soft: true,
		})
	}

	if exit, ok := err.(Exit); ok == true {
		panic(exit)
	}

	panic(Exit{
		Code: code,
		Msg:  err.Error(),
		Soft: true,
	})
}

func ExitOrDieFromError(code int, err error, isSoftConditions ...bool) {
	if err == nil {
		panic(Exit{
			Code: code,
			Soft: true,
		})
	}

	if exit, ok := err.(Exit); ok == true {
		panic(exit)
	}

	soft := zerrors.IsErrorSoft(err)

	// find any soft condition
	for idx := 0; !soft && idx < len(isSoftConditions); idx++ {
		soft = isSoftConditions[idx]
	}

	panic(Exit{
		Code: code,
		Msg:  err.Error(),
		Soft: soft,
	})
}

// logic should equal to utilstest.TestHandleExit func
func HandleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true && exit.Soft {
			errMsg := exit.Error()
			if len(errMsg) > 0 {
				if errMsg[len(errMsg)-1] != '\n' {
					errMsg += "\n"
				}
				_, _ = os.Stderr.WriteString(errMsg)
			}
			os.Exit(exit.Code)
		}
		panic(e)
	}
}
