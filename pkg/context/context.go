package context

import (
	"sort"

	"chaos-mesh/matrix/pkg/node"
	"chaos-mesh/matrix/pkg/random"
	"chaos-mesh/matrix/pkg/serializer"
	"chaos-mesh/matrix/pkg/synthesizer"
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
func (c MatrixContext) Gen(seed int64) node.ConfigGroup {
	var result node.ConfigGroup

	result.Configs = make(map[serializer.Config]interface{})

	random.Seed(seed)
	for _, k := range sortedKeys(c.Configs) {
		result.Configs[c.Configs[k].Config] = synthesizer.SimpleRecGen(c.Configs[k].Hollow)
	}

	return result
}
