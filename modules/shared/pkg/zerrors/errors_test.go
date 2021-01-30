package zerrors_test

import (
	"errors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	httperrors "k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

var _ = Describe("Errors", func() {
	table.DescribeTable("should distinguish soft errors", func(tested error, result bool) {
		Expect(zerrors.IsErrorSoft(tested)).To(Equal(result))
	},
		table.Entry("nil", nil, false),
		table.Entry("not soft", errors.New("not soft"), false),
		table.Entry("not soft", zerrors.NewMissingRequiredError("soft"), true),
	)

	table.DescribeTable("should detect status errors to be soft", func(tested error, allowedCodes []int32, result bool) {
		Expect(zerrors.IsStatusError(tested, allowedCodes...)).To(Equal(result))
	},
		table.Entry("nil", nil, nil, false),
		table.Entry("not status error", errors.New("not soft"), nil, false),
		table.Entry("not status error (soft)", zerrors.NewMissingRequiredError("soft"), nil, false),
		table.Entry("status error with no allowed codes", httperrors.NewUnauthorized("unauthorized"), nil, false),
		table.Entry("status error with wrong allowed codes", httperrors.NewUnauthorized("unauthorized"), []int32{http.StatusConflict, http.StatusNotFound}, false),
		table.Entry("status error with correct allowed codes", httperrors.NewUnauthorized("unauthorized"), []int32{http.StatusUnauthorized, http.StatusNotFound}, true),
	)
})
