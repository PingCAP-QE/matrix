// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"testing"

	"github.com/chaos-mesh/matrix/pkg/node/data"

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
