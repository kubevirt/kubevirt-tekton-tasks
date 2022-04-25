package zerrors_test

import (
	"errors"
	"net/http"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	httperrors "k8s.io/apimachinery/pkg/api/errors"
)

var _ = Describe("Errors", func() {
	DescribeTable("should distinguish soft errors", func(tested error, result bool) {
		Expect(zerrors.IsErrorSoft(tested)).To(Equal(result))
	},
		Entry("nil", nil, false),
		Entry("not soft", errors.New("not soft"), false),
		Entry("not soft", zerrors.NewMissingRequiredError("soft"), true),
	)

	DescribeTable("should detect status errors to be soft", func(tested error, allowedCodes []int32, result bool) {
		Expect(zerrors.IsStatusError(tested, allowedCodes...)).To(Equal(result))
	},
		Entry("nil", nil, nil, false),
		Entry("not status error", errors.New("not soft"), nil, false),
		Entry("not status error (soft)", zerrors.NewMissingRequiredError("soft"), nil, false),
		Entry("status error with no allowed codes", httperrors.NewUnauthorized("unauthorized"), nil, false),
		Entry("status error with wrong allowed codes", httperrors.NewUnauthorized("unauthorized"), []int32{http.StatusConflict, http.StatusNotFound}, false),
		Entry("status error with correct allowed codes", httperrors.NewUnauthorized("unauthorized"), []int32{http.StatusUnauthorized, http.StatusNotFound}, true),
	)
})
