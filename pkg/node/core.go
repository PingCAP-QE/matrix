package node

import (
	"chaos-mesh/matrix/pkg/node/data"
	"chaos-mesh/matrix/pkg/serializer"
)

type ConfigGroup struct {
	Configs map[serializer.Config]data.ConcreteInterface
}
