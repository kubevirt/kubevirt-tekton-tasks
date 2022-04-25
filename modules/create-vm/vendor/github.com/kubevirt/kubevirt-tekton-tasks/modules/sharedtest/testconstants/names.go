package testconstants

import (
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/onsi/ginkgo/v2"
)

func fiveDigitTestHash(s string) string {
	bitCount := 5 * 4
	clearBitCount := 32 - bitCount
	h := fnv.New64()
	_, _ = h.Write([]byte(s))

	hash := h.Sum64()

	// mix with the seed and node number
	hash = hash ^ uint64(ginkgo.GinkgoRandomSeed()+int64(ginkgo.GinkgoParallelNode()))

	// xor lower 32 bits with higher 32 bits and keep only lower bits
	hash32 := uint32(((hash ^ (hash >> 32)) << 32) >> 32)

	// xor highest 12 bits with lower ones and keep only lower bits
	hash32 = ((hash32 ^ (hash32 >> bitCount)) << clearBitCount) >> clearBitCount

	// forcefully add left fifth digit if lower number
	lowest5DigitNum := uint32(1) << (bitCount - 4)
	if hash32 < lowest5DigitNum {
		hash32 ^= lowest5DigitNum
	}
	// will result in 5 places
	return strconv.FormatUint(uint64(hash32), 16)
}

func TestRandomName(name string) string {
	// convert Full Test description into ID
	id := fiveDigitTestHash(ginkgo.CurrentGinkgoTestDescription().FullTestText)

	return strings.Join([]string{name, id}, "-")
}
