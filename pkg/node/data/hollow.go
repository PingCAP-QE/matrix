package data

type HollowInterface interface {
	HollowType() string
}

type ValuedHollowInterface interface {
	HollowValue() interface{}
}

type HollowInt struct {
	Range []int
}

func (h HollowInt) HollowType() string { return TypeInt }

type HollowBool struct{}

func (h HollowBool) HollowType() string { return TypeBool }

type HollowString struct{}

func (h HollowString) HollowType() string { return TypeString }

type HollowMap struct {
	Map map[string]interface{}
}

func (h HollowMap) HollowType() string       { return TypeMap }
func (h HollowMap) HollowValue() interface{} { panic("not implemented") }

var _ HollowInterface = (*HollowInt)(nil)
var _ HollowInterface = (*HollowBool)(nil)
var _ HollowInterface = (*HollowString)(nil)
var _ HollowInterface = (*HollowMap)(nil)
var _ ValuedHollowInterface = (*HollowMap)(nil)

func NewHollowInt(start int, end int) HollowInt {
	h := HollowInt{}
	h.Range = make([]int, 2)
	h.Range[0] = start
	h.Range[1] = end

	return h
}
