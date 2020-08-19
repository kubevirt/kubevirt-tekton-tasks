package errors_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/errors"
	"strconv"
)

var _ = Describe("MultiError", func() {
	Describe("test CRUD", func() {
		table.DescribeTable("should be empty by default", func(err *errors2.MultiError) {
			Expect(err.IsEmpty()).To(BeTrue())
			Expect(err.Len()).To(BeZero())
			Expect(err.AsOptional()).To(BeNil())
		},
			table.Entry("nil", nil),
			table.Entry("new", &errors2.MultiError{}),
			table.Entry("new func", errors2.NewMultiError()),
		)

		It("add and get works and error is not optional", func() {
			err := errors2.NewMultiError()

			for i := 0; i < 7; i++ {
				id := strconv.Itoa(i)
				msg := "msg " + id
				if i > 4 {
					err = err.AddC(id, errors.New(msg)) // try both versions
				} else {
					err.Add(id, errors.New(msg))
				}

				Expect(err.IsEmpty()).To(BeFalse())
				Expect(err.Len()).To(Equal(i + 1))
				Expect(err.Get(id).Error()).To(Equal(msg))
				Expect(errors2.GetErrorFromMultiError(err, id).Error()).To(Equal(msg))
				Expect(err.AsOptional()).To(Equal(err))
			}
		})
	})

	Describe("prints correct messages", func() {
		err := errors2.NewMultiError().
			AddC("1", errors.New("err1")).
			AddC("2", errors.New("err2")).
			AddC("3", errors.New("err3"))

		It("long message", func() {
			longResults := "err1\nerr2\nerr3\n"
			Expect(err.Error()).To(Equal(longResults))
			Expect(err.LongPrint().Error()).To(Equal(longResults))
		})

		It("short message", func() {
			shortResults := "errs: 1, 2, 3"

			Expect(err.ShortPrint("errs:").Error()).To(Equal(shortResults))
		})

		It("empty message", func() {
			Expect(errors2.NewMultiError().Error()).To(BeEmpty())
		})
	})

	table.DescribeTable("correctly reports soft errors", func(tested *errors2.MultiError, result bool) {
		Expect(tested.IsSoft()).To(Equal(result))
		Expect(errors2.IsErrorSoft(tested)).To(Equal(result))
	},
		table.Entry("nil soft", nil, true),
		table.Entry("new soft", &errors2.MultiError{}, true),
		table.Entry("new func soft", errors2.NewMultiError(), true),
		table.Entry("soft with only soft errors", errors2.NewMultiError().
			AddC("soft1", errors2.NewMissingRequiredError("soft1")).
			AddC("soft2", errors2.NewMissingRequiredError("soft2")), true),
		table.Entry("not soft with one hard and one soft", errors2.NewMultiError().
			AddC("soft1", errors2.NewMissingRequiredError("soft1")).
			AddC("hard2", errors.New("hard2")), false),
		table.Entry("not soft with only hard", errors2.NewMultiError().
			AddC("hard2", errors.New("hard2")).
			AddC("hard2", errors.New("hard2")), false),
	)

})
