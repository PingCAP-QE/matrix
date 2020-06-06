package context

import (
	"chaos-mesh/matrix/pkg/node"
)

type MatrixContext struct {
	Configs map[string]node.AbstractConfig
}

// This is to generate real value from an abstract tree
func (MatrixContext) Gen() node.ConfigGroup {
	panic("not implemented")
}
