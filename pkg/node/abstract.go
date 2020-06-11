package node

import (
	"chaos-mesh/matrix/pkg/serializer"
)

type AbstractConfig struct {
	Tag string
	serializer.Config
	Hollow interface{}
}
