package constants

import (
	"github.com/onsi/ginkgo"
	"hash/fnv"
	"strconv"
	"strings"
	"time"
)

const e2eNamespacePrefix = "e2e-tests"

const (
	PollInterval = 1 * time.Second
)

type TargetNamespace string

const (
	DeployTargetNS TargetNamespace = "deploy"
	TestTargetNS   TargetNamespace = "test"
	CustomTargetNS TargetNamespace = "custom"
)

func E2ETestsName(name string) string {
	// convert Full Test description into ID
	id := fiveDigitTestHash(ginkgo.CurrentGinkgoTestDescription().FullTestText)

	return strings.Join([]string{e2eNamespacePrefix, name, id}, "-")
}

func ToStringBoolean(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func fiveDigitTestHash(s string) string {
	bitCount := 5 * 4
	clearBitCount := 32 - bitCount
	h := fnv.New64()
	_, _ = h.Write([]byte(s))

	hash := h.Sum64()

	// mix with the seed
	hash = hash ^ uint64(ginkgo.GinkgoRandomSeed())

	// xor lower 32 bits with higher 32 bits and keep only lower bits
	hash32 := uint32(((hash ^ (hash >> 32)) << 32) >> 32)

	// xor highest 12 bits with lower ones and keep only lower bits
	hash32 = ((hash32 ^ (hash32 >> bitCount)) << clearBitCount) >> clearBitCount

	// will result in 5 places
	return strconv.FormatUint(uint64(hash32), 16)
}
