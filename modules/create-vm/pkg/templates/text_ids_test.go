package templates

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"math/rand"
	"sort"
)

var _ = Describe("TextIds", func() {

	table.DescribeTable("check correct ascending order", func(ids []string) {
		toSort := textIDs{}
		toSort = append(toSort, ids...)
		result := textIDs{}
		result = append(result, ids...)

		rand.Seed(GinkgoRandomSeed())
		rand.Shuffle(len(toSort), func(i, j int) { toSort[i], toSort[j] = toSort[j], toSort[i] })

		sort.Sort(toSort)

		Expect(toSort).To(Equal(result))
	},
		table.Entry("ids", []string{
			"os-a-0.",
			"os-a-0.0,0",
			"os-a-1",
			"os-a-1.0",
			"os-a-1.0,0",
			"os-a-1.1",
			"os-a-1.2,3",
			"os-a-1.2,4",
			"os-a-2",
		}),
		table.Entry("fedora", []string{
			"fedora27",
			"silverblue28",
			"fedora28",
			"silverblue29",
			"fedora29",
		}),
		table.Entry("win", []string{
			"win2k8",
			"win2k8r2",
			"win2k12r2",
			"win2k16",
			"win2k19",
			"win10",
		}),
		table.Entry("win", []string{
			"win2k8",
			"win2k8r2",
			"win2k12r2",
			"win2k16",
			"win2k19",
			"win10",
		}),
		table.Entry("rhel", []string{
			"rhel7.0",
			"rhel7.1",
			"rhel7.2",
			"rhel7.3",
			"rhel7.10",
			"rhel7.11",
			"rhel8.1",
			"rhel8.2",
		}),
		table.Entry("text based", []string{
			"centos",
			"fedora",
			"rhel",
			"ubuntu",
			"win",
		}),
	)

})
