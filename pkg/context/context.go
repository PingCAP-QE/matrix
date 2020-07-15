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

package context

import (
	"errors"
	"fmt"
	"path"
	"sort"

	"github.com/chaos-mesh/matrix/pkg/node"
	"github.com/chaos-mesh/matrix/pkg/random"
	"github.com/chaos-mesh/matrix/pkg/serializer"
	"github.com/chaos-mesh/matrix/pkg/synthesizer"
)

type MatrixContext struct {
	Configs map[string]node.AbstractConfig
}

func sortedKeys(configs map[string]node.AbstractConfig) []string {
	var sortedKeys []string
	for i, _ := range configs {
		sortedKeys = append(sortedKeys, i)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// This is to generate real value from an abstract tree
func (c MatrixContext) gen(seed int64) node.ConfigGroup {
	var result node.ConfigGroup

	result.Configs = make(map[serializer.Config]interface{})

	random.Seed(seed)
	for _, k := range sortedKeys(c.Configs) {
		result.Configs[c.Configs[k].Config] = synthesizer.SimpleRecGen(c.Configs[k].Hollow)
	}

	return result
}

func (c MatrixContext) Dump(seed int64, output string) error {
	var err error
	values := c.gen(seed)
	for config, concrete := range values.Configs {
		err = config.Serializer.Dump(concrete, path.Join(output, config.Target))
		if err != nil {
			return errors.New(fmt.Sprintf("error %s when dumping %v", err.Error(), concrete))
		}
	}
	return nil
}
