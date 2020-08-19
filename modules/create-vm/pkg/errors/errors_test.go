package errors_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/errors"
	errors3 "k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

var _ = Describe("Errors", func() {
	table.DescribeTable("should distinguish soft errors", func(tested error, result bool) {
		Expect(errors2.IsErrorSoft(tested)).To(Equal(result))
	},
		table.Entry("nil", nil, false),
		table.Entry("not soft", errors.New("not soft"), false),
		table.Entry("not soft", errors2.NewMissingRequiredError("soft"), true),
	)

	table.DescribeTable("should detect status errors to be soft", func(tested error, allowedCodes []int32, result bool) {
		Expect(errors2.IsStatusErrorSoft(tested, allowedCodes...)).To(Equal(result))
	},
		table.Entry("nil", nil, nil, false),
		table.Entry("not status error", errors.New("not soft"), nil, false),
		table.Entry("not status error (soft)", errors2.NewMissingRequiredError("soft"), nil, false),
		table.Entry("status error with no allowed codes", errors3.NewUnauthorized("unauthorized"), nil, false),
		table.Entry("status error with wrong allowed codes", errors3.NewUnauthorized("unauthorized"), []int32{http.StatusConflict, http.StatusNotFound}, false),
		table.Entry("status error with correct allowed codes", errors3.NewUnauthorized("unauthorized"), []int32{http.StatusUnauthorized, http.StatusNotFound}, true),
	)
})
