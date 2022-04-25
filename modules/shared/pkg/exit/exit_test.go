package exit_test

import (
	"errors"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	. "github.com/onsi/ginkgo/v2"
)

const (
	exFailed        = "execution failed"
	exFailedRes     = exFailed + "\n"
	exFailedLong    = "execution failed\n additional info\n"
	exFailedLongRes = exFailedLong
)

var _ = Describe("Utils", func() {
	Describe("exits correctly", func() {
		DescribeTable("just exit", func(err error, shouldExitWithCode int, shouldExitWithMessage string) {
			defer utilstest.HandleTestExit(false, shouldExitWithCode, shouldExitWithMessage)
			exit.ExitFromError(shouldExitWithCode, err)
		},
			Entry("no err", nil, 0, ""),
			Entry("err", errors.New(exFailed), 2, exFailedRes),
			Entry("long err", errors.New(exFailedLong), 3, exFailedLongRes),
		)
	})

	Describe("exits or dies correctly", func() {
		DescribeTable("exit or die", func(err error, shouldExitWithCode int, shouldExitWithMessage string, shouldPanic bool, softConditions []bool) {
			defer utilstest.HandleTestExit(shouldPanic, shouldExitWithCode, shouldExitWithMessage)
			exit.ExitOrDieFromError(shouldExitWithCode, err, softConditions...)
		},
			Entry("no err exits", nil, 0, "", false, nil),
			Entry("hard err dies", errors.New(exFailed), 2, exFailedRes, true, []bool{false}),
			Entry("soft err exits", zerrors.NewMissingRequiredError(exFailed), 3, exFailedRes, false, nil),
			Entry("hard err with additional soft conditions exits", errors.New(exFailed), 2, exFailedRes, false, []bool{false, true, false}),
		)
	})
})
