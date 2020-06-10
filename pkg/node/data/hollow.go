package data

import (
	"time"

	"github.com/c2h5oh/datasize"
)

type HollowInterface interface {
	GetType() string
	GetAllCondition() ([]HollowCondition, bool)
}

type HollowCondition struct{ Raw []interface{} }

type HollowBool struct{ Condition *HollowCondition }

type HollowInt struct {
	Range     []int
	Condition *HollowCondition
}
type HollowFloat struct {
	Range     []float64
	Condition *HollowCondition
}
type HollowString struct {
	Value     string
	Condition *HollowCondition
}
type HollowTime struct {
	Range     []time.Duration
	Condition *HollowCondition
}
type HollowSize struct {
	Range     []datasize.ByteSize
	Condition *HollowCondition
}
type HollowChoice struct {
	List      []interface{}
	Condition *HollowCondition
}
type HollowList struct {
	List      []interface{}
	Condition *HollowCondition
}
type HollowMap struct {
	Map       map[string]interface{}
	Condition *HollowCondition
}

var _ HollowInterface = (*HollowBool)(nil)
var _ HollowInterface = (*HollowInt)(nil)
var _ HollowInterface = (*HollowFloat)(nil)
var _ HollowInterface = (*HollowString)(nil)
var _ HollowInterface = (*HollowTime)(nil)
var _ HollowInterface = (*HollowSize)(nil)
var _ HollowInterface = (*HollowChoice)(nil)
var _ HollowInterface = (*HollowList)(nil)
var _ HollowInterface = (*HollowMap)(nil)

func (h HollowBool) GetType() string   { return TypeBool }
func (h HollowInt) GetType() string    { return TypeInt }
func (h HollowFloat) GetType() string  { return TypeFloat }
func (h HollowString) GetType() string { return TypeString }
func (h HollowTime) GetType() string   { return TypeTime }
func (h HollowSize) GetType() string   { return TypeSize }
func (h HollowChoice) GetType() string { return TypeChoice }
func (h HollowList) GetType() string   { return TypeList }
func (h HollowMap) GetType() string    { return TypeMap }

func returnCond(condition *HollowCondition) ([]HollowCondition, bool) {
	if condition == nil {
		return []HollowCondition{}, false
	}
	return []HollowCondition{*condition}, true
}

func returnListCond(condition *HollowCondition, list []interface{}) ([]HollowCondition, bool) {
	condList, _ := returnCond(condition)
	for _, v := range list {
		childCond, _ := v.(HollowInterface).GetAllCondition()
		for _, cond := range childCond {
			condList = append(condList, cond)
		}
	}
	return condList, len(condList) > 0
}

func (h HollowBool) GetAllCondition() ([]HollowCondition, bool)   { return returnCond(h.Condition) }
func (h HollowInt) GetAllCondition() ([]HollowCondition, bool)    { return returnCond(h.Condition) }
func (h HollowFloat) GetAllCondition() ([]HollowCondition, bool)  { return returnCond(h.Condition) }
func (h HollowString) GetAllCondition() ([]HollowCondition, bool) { return returnCond(h.Condition) }
func (h HollowTime) GetAllCondition() ([]HollowCondition, bool)   { return returnCond(h.Condition) }
func (h HollowSize) GetAllCondition() ([]HollowCondition, bool)   { return returnCond(h.Condition) }
func (h HollowChoice) GetAllCondition() ([]HollowCondition, bool) {
	return returnListCond(h.Condition, h.List)
}
func (h HollowList) GetAllCondition() ([]HollowCondition, bool) {
	return returnListCond(h.Condition, h.List)
}
func (h HollowMap) GetAllCondition() ([]HollowCondition, bool) {
	condList, _ := returnCond(h.Condition)
	for _, v := range h.Map {
		childCond, _ := v.(HollowInterface).GetAllCondition()
		for _, cond := range childCond {
			condList = append(condList, cond)
		}
	}
	return condList, len(condList) > 0
}

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
