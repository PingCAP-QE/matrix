package data

type HollowInterface interface {
	HollowType() string
}

type ValuedHollowInterface interface {
	HollowValue() interface{}
}

type HollowInt struct{}

func (h HollowInt) HollowType() string { return TypeInt }

type HollowMap struct{}

func (h HollowMap) HollowType() string       { return TypeMap }
func (h HollowMap) HollowValue() interface{} { panic("not implemented") }
