// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"github.com/c2h5oh/datasize"
)

type HollowInterface interface {
	GetType() string
	GetAllCondition() ([]HollowCondition, bool)
	SetCondition(condition *HollowCondition)
}

type HollowCondition struct{ Raw []interface{} }

type HollowBool struct {
	Value     *bool
	Condition *HollowCondition
}

type HollowInt struct {
	RangeStart, RangeEnd int
	Condition            *HollowCondition
}
type HollowFloat struct {
	RangeStart, RangeEnd float64
	Condition            *HollowCondition
}
type HollowString struct {
	Value     string
	Condition *HollowCondition
}
type HollowTime struct {
	RangeStart, RangeEnd Time
	Condition            *HollowCondition
}
type HollowSize struct {
	RangeStart, RangeEnd datasize.ByteSize
	Condition            *HollowCondition
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

func (h HollowBool) SetCondition(condition *HollowCondition)   { h.Condition = condition }
func (h HollowInt) SetCondition(condition *HollowCondition)    { h.Condition = condition }
func (h HollowFloat) SetCondition(condition *HollowCondition)  { h.Condition = condition }
func (h HollowString) SetCondition(condition *HollowCondition) { h.Condition = condition }
func (h HollowTime) SetCondition(condition *HollowCondition)   { h.Condition = condition }
func (h HollowSize) SetCondition(condition *HollowCondition)   { h.Condition = condition }
func (h HollowChoice) SetCondition(condition *HollowCondition) { h.Condition = condition }
func (h HollowList) SetCondition(condition *HollowCondition)   { h.Condition = condition }
func (h HollowMap) SetCondition(condition *HollowCondition)    { h.Condition = condition }
