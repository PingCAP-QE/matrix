package parser

import (
	"testing"

	"chaos-mesh/matrix/pkg/node/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Parser Suite")
}

var _ = Describe("Parser", func() {
	Context("Parser", func() {
		It("Parser Types", func() {
			initValueParserMap()
			for _, dataType := range data.AllTypes {
				_, ok := valueParserMap[dataType]
				Expect(ok).To(BeTrue())
			}
		})
	})
})
