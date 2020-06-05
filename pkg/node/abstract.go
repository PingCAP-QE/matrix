package node

import "chaos-mesh/matrix/pkg/serializer"

type HollowInterface interface {
	HollowValue() interface{}
}

var _ HollowInterface = (*Hollow)(nil)

type Hollow struct {
	Value interface{}
}

func (h Hollow) HollowValue() interface{} {
	return h.Value
}

type AbstractConfig struct {
	Tag string
	serializer.Config
	Hollow
}
