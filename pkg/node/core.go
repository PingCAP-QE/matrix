package node

import (
	"chaos-mesh/matrix/pkg/serializer"
)

type ConfigGroup struct {
	Configs map[serializer.Config]interface{}
}
