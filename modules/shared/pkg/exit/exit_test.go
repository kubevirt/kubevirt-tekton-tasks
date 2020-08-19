package exit_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/utilstest"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
)

const (
	exFailed        = "execution failed"
	exFailedRes     = exFailed + "\n"
	exFailedLong    = "execution failed\n additional info\n"
	exFailedLongRes = exFailedLong
)

var _ = Describe("Utils", func() {
	Describe("exits correctly", func() {
		table.DescribeTable("just exit", func(err error, shouldExitWithCode int, shouldExitWithMessage string) {
			defer utilstest.HandleTestExit(false, shouldExitWithCode, shouldExitWithMessage)
			exit.ExitFromError(shouldExitWithCode, err)
		},
			table.Entry("no err", nil, 0, ""),
			table.Entry("err", errors.New(exFailed), 2, exFailedRes),
			table.Entry("long err", errors.New(exFailedLong), 3, exFailedLongRes),
		)
	})

	Describe("exits or dies correctly", func() {
		table.DescribeTable("exit or die", func(err error, shouldExitWithCode int, shouldExitWithMessage string, shouldPanic bool, softConditions []bool) {
			defer utilstest.HandleTestExit(shouldPanic, shouldExitWithCode, shouldExitWithMessage)
			exit.ExitOrDieFromError(shouldExitWithCode, err, softConditions...)
		},
			table.Entry("no err exits", nil, 0, "", false, nil),
			table.Entry("hard err dies", errors.New(exFailed), 2, exFailedRes, true, []bool{false}),
			table.Entry("soft err exits", zerrors.NewMissingRequiredError(exFailed), 3, exFailedRes, false, nil),
			table.Entry("hard err with additional soft conditions exits", errors.New(exFailed), 2, exFailedRes, false, []bool{false, true, false}),
		)
	})
})
