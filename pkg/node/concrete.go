package node

import "chaos-mesh/matrix/pkg/serializer"

// Concrete value
type Concrete interface {
	Value() interface{}
}

type ConfigGroup struct {
	Configs map[serializer.Config]Concrete
}
