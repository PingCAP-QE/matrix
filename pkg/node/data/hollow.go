package data

type HollowInterface interface {
	HollowType() string
}

type HollowBranch struct {
	Value interface{}
	When  interface{}
}

type HollowBool struct{}
type HollowInt struct{ Range []int }
type HollowFloat struct{ Range []float64 }
type HollowString struct{ Value string }
type HollowChoice struct{ List []HollowBranch }
type HollowList struct{ List []interface{} }
type HollowMap struct{ Map map[string]interface{} }

var _ HollowInterface = (*HollowBool)(nil)
var _ HollowInterface = (*HollowInt)(nil)
var _ HollowInterface = (*HollowFloat)(nil)
var _ HollowInterface = (*HollowString)(nil)
var _ HollowInterface = (*HollowBranch)(nil)
var _ HollowInterface = (*HollowChoice)(nil)
var _ HollowInterface = (*HollowList)(nil)
var _ HollowInterface = (*HollowMap)(nil)

func (h HollowBool) HollowType() string   { return TypeBool }
func (h HollowInt) HollowType() string    { return TypeInt }
func (h HollowFloat) HollowType() string  { return TypeFloat }
func (h HollowString) HollowType() string { return TypeString }
func (h HollowBranch) HollowType() string { return TypeBranch }
func (h HollowChoice) HollowType() string { return TypeChoice }
func (h HollowList) HollowType() string   { return TypeList }
func (h HollowMap) HollowType() string    { return TypeMap }

func NewHollowInt(start int, end int) HollowInt {
	h := HollowInt{}
	h.Range = make([]int, 2)
	h.Range[0] = start
	h.Range[1] = end
	return h
}

func NewHollowFloat(start float64, end float64) HollowFloat {
	h := HollowFloat{}
	h.Range = make([]float64, 2)
	h.Range[0] = start
	h.Range[1] = end
	return h
}
