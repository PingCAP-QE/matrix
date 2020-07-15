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

package parser

import (
	"errors"
	"fmt"

	"github.com/c2h5oh/datasize"

	"github.com/pingcap/log"

	"github.com/chaos-mesh/matrix/pkg/node/data"
	"github.com/chaos-mesh/matrix/pkg/utils"
)

const ExprNotSupportedMessage = "expr not supported"

func parseTree(rawValue interface{}) (interface{}, error) {
	switch rawValue.(type) {
	case map[string]interface{}:
		return parseMap(rawValue.(map[string]interface{}))
	case []interface{}:
		return parseChoice(rawValue.([]interface{}))
	default:
		return parseLiteral(rawValue)
	}
}

func parseLiteral(rawLiteral interface{}) (interface{}, error) {
	switch rawLiteral.(type) {
	case float64:
		return parseFloat(rawLiteral.(float64)), nil
	case bool:
		b := rawLiteral.(bool)
		return data.HollowBool{Value: &b}, nil
	case string:
		return parseString(rawLiteral.(string))
	default:
		log.L().Warn(fmt.Sprintf("%s not handled, return HollowInt", rawLiteral))
		return data.HollowInt{RangeStart: utils.MinInt, RangeEnd: utils.MaxInt}, nil
	}
}

func parseValueMap(rawMap map[string]interface{}) (interface{}, error) {
	valueType, ok := rawMap["type"]
	if !ok {
		return nil, errors.New("no `type` field")
	}
	stringValueType, ok := valueType.(string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("`type` of %v is not of type string", rawMap))
	}
	if valueParserMap == nil {
		initValueParserMap()
	}
	valueParser, ok := valueParserMap[stringValueType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no parser for type \"%s\"", stringValueType))
	}
	hollowValue, err := valueParser(rawMap)
	if err == nil {
		if rawCondition, ok := rawMap["when"]; ok {
			delete(rawMap, "when")
			var hollowCondition *data.HollowCondition
			hollowCondition, err = parseCondition(rawCondition)
			if err != nil {
				return nil, err
			}
			hollowValue.(data.HollowInterface).SetCondition(hollowCondition)
		}

		if _, ok = rawMap["default"]; ok {
			// todo: store default value
			delete(rawMap, "default")
		}
	}
	return hollowValue, err
}

func parseMap(rawMap map[string]interface{}) (interface{}, error) {
	var err error
	hollowValue, err := parseValueMap(rawMap)
	if err == nil {
		return hollowValue, err
	}
	var hollowMap data.HollowMap
	hollowMap.Map = make(map[string]interface{})
	for k := range rawMap {
		var v interface{}
		v, err = parseTree(rawMap[k])
		if err != nil {
			return nil, err
		}
		hollowMap.Map[k] = v
	}
	return hollowMap, nil
}

func parseChoice(rawList []interface{}) (interface{}, error) {
	var hollowChoice data.HollowChoice
	if len(rawList) == 0 {
		return nil, errors.New("no available option for choice")
	}
	hollowChoice.List = make([]interface{}, len(rawList))
	for i, k := range rawList {
		hollowBranch, err := parseTree(k)
		if err != nil {
			return nil, err
		}
		hollowChoice.List[i] = hollowBranch
	}
	return hollowChoice, nil
}

func parseCondition(rawCond interface{}) (*data.HollowCondition, error) {
	switch rawCond.(type) {
	case string:
	case []interface{}:
	default:
		log.L().Warn(fmt.Sprintf("ignoring constraint: %s", rawCond))
	}
	return nil, nil
}

func parseDefault(rawString string) (interface{}, error) {
	if hollow, ok := data.DefaultValue[rawString]; ok {
		return hollow, nil
	}
	return nil, errors.New(fmt.Sprintf("not default value for %s", rawString))
}

// parseString convert string literals to corresponding values, default value will be handled here
func parseString(rawString string) (interface{}, error) {
	defaultValue, err := parseDefault(rawString)
	if err == nil {
		return defaultValue, nil
	} else {
		return data.HollowString{Value: rawString}, nil
	}
}

func parseFloat(rawFloat float64) interface{} {
	// int / float literal
	intValue := int(rawFloat)
	if rawFloat == float64(intValue) {
		return data.HollowInt{RangeStart: intValue, RangeEnd: intValue}
	} else {
		return data.HollowFloat{RangeStart: rawFloat, RangeEnd: rawFloat}
	}
}

func parseHollowInt(raw map[string]interface{}) (interface{}, error) {
	var hollowInt data.HollowInt
	var err error

	if intRange, ok := raw["range"]; ok {
		hollowInt.RangeStart, hollowInt.RangeEnd, err = parseIntRange(intRange)
		if err != nil {
			return nil, err
		}
	}
	return hollowInt, nil
}

func parseHollowFloat(raw map[string]interface{}) (interface{}, error) {
	var hollowFloat data.HollowFloat
	var err error

	if floatRange, ok := raw["range"]; ok {
		hollowFloat.RangeStart, hollowFloat.RangeEnd, err = parseFloatRange(floatRange)
		if err != nil {
			return nil, err
		}
	}
	return hollowFloat, nil
}

func parseHollowTime(raw map[string]interface{}) (interface{}, error) {
	var hollowTime data.HollowTime
	var err error

	if timeRange, ok := raw["range"]; ok {
		hollowTime.RangeStart, hollowTime.RangeEnd, err = parseTimeRange(timeRange)
		if err != nil {
			return nil, err
		}
	}
	return hollowTime, nil
}

func parseHollowSize(raw map[string]interface{}) (interface{}, error) {
	var hollowSize data.HollowSize
	var err error

	if timeRange, ok := raw["range"]; ok {
		hollowSize.RangeStart, hollowSize.RangeEnd, err = parseSizeRange(timeRange)
		if err != nil {
			return nil, err
		}
	}
	return hollowSize, nil
}

func parseIntRange(raw interface{}) (int, int, error) {
	rangeList := raw.([]interface{})
	dur := make([]int, 2)
	for i, v := range rangeList {
		switch v.(type) {
		case float64:
			dur[i] = int(v.(float64))
		default:
			return 0, 0, errors.New(fmt.Sprintf("%s cannot be parsed as int", v))
		}
	}
	return dur[0], dur[1], nil
}

func parseFloatRange(raw interface{}) (float64, float64, error) {
	rangeList := raw.([]interface{})
	dur := make([]float64, 2)
	for i, v := range rangeList {
		switch v.(type) {
		case float64:
			dur[i] = v.(float64)
		default:
			return 0, 0, errors.New(fmt.Sprintf("%s cannot be parsed as float", v))
		}
	}
	return dur[0], dur[1], nil
}

func parseTimeRange(raw interface{}) (data.Time, data.Time, error) {
	rangeList := raw.([]interface{})
	dur := make([]data.Time, 2)
	for i, v := range rangeList {
		if timeStr, ok := v.(string); ok {
			t, err := data.NewTimeFromString(timeStr)
			if err != nil {
				return data.NewTime(0), data.NewTime(0), err
			}
			dur[i] = *t
		} else {
			return data.NewTime(0), data.NewTime(0), errors.New(fmt.Sprintf("%s cannot be parsed as time", v))
		}
	}
	return dur[0], dur[1], nil
}

func parseSizeRange(raw interface{}) (datasize.ByteSize, datasize.ByteSize, error) {
	var err error
	rangeList := raw.([]interface{})
	dur := make([]datasize.ByteSize, 2)
	for i, v := range rangeList {
		if sizeStr, ok := v.(string); ok {
			var size datasize.ByteSize
			err = size.UnmarshalText([]byte(sizeStr))
			if err != nil {
				return datasize.ByteSize(0), datasize.ByteSize(0), err
			}
			dur[i] = size
		} else {
			return datasize.ByteSize(0), datasize.ByteSize(0), errors.New(fmt.Sprintf("%s cannot be parsed as size", v))
		}
	}
	return dur[0], dur[1], nil
}

func parseHollowList(raw map[string]interface{}) (interface{}, error) {
	if rawList, ok := raw["value"]; ok {
		if rawList, ok := rawList.([]interface{}); ok {
			var list data.HollowList
			var err error
			list.List = make([]interface{}, len(rawList))
			for i, v := range rawList {
				list.List[i], err = parseTree(v)
				if err != nil {
					return nil, err
				}
			}
			return list, nil
		} else {
			return nil, errors.New(fmt.Sprintf("field `value` of type `list` is not if type `[]interface{}`: %s", rawList))
		}
	} else {
		return nil, errors.New("type `list` does not contain field `value`")
	}
}

func parseHollowChoice(raw map[string]interface{}) (interface{}, error) {
	if rawList, ok := raw["value"]; ok {
		if rawList, ok := rawList.([]interface{}); ok {
			return parseChoice(rawList)
		} else {
			return nil, errors.New(fmt.Sprintf("field `value` of type `choice` is not if type `[]interface{}`: %s", rawList))
		}
	} else {
		return nil, errors.New("type `choice` does not contain field `value`")
	}
}

func parseHollowChoiceN(raw map[string]interface{}) (interface{}, error) {
	hollowChoice, err := parseHollowChoice(raw)
	if err != nil {
		return nil, err
	}
	hollowChoiceN := data.HollowChoiceN{HollowChoice: hollowChoice.(data.HollowChoice)}
	if n, ok := raw["n"]; ok {
		if intN, ok := n.(int); ok {
			hollowChoiceN.N = intN
		} else {
			return nil, errors.New(fmt.Sprintf("`n` of choice_n is not of type int: %v", n))
		}
	}
	if sep, ok := raw["sep"]; ok {
		if strSep, ok := sep.(string); ok {
			hollowChoiceN.Sep = strSep
		} else {
			return nil, errors.New(fmt.Sprintf("`sep` of choice_n is not of type string: %v", sep))
		}
	}
	return hollowChoiceN, nil
}

func parseHollowMap(raw map[string]interface{}) (interface{}, error) {
	if rawList, ok := raw["value"]; ok {
		if rawMap, ok := rawList.(map[string]interface{}); ok {
			var hollowMap data.HollowMap
			var err error
			hollowMap.Map = make(map[string]interface{}, len(rawMap))
			for i, v := range rawMap {
				hollowMap.Map[i], err = parseTree(v)
				if err != nil {
					return nil, err
				}
			}
			return hollowMap, nil
		} else {
			return nil, errors.New(fmt.Sprintf("field `value` of type `map` is not if type `map[string]interface{}`: %s", rawMap))
		}
	} else {
		return nil, errors.New("type `map` does not contain field `value`")
	}
}

var valueParserMap map[string]func(map[string]interface{}) (interface{}, error)

func initValueParserMap() {
	valueParserMap = map[string]func(map[string]interface{}) (interface{}, error){
		data.TypeBool: func(rawBool map[string]interface{}) (interface{}, error) {
			if rawBoolValue, ok := rawBool["value"]; ok {
				if b, ok := rawBoolValue.(bool); ok {
					return data.HollowBool{Value: &b}, nil
				} else {
					return nil, errors.New("`value` field of bool is not of type bool")
				}
			}
			return data.HollowBool{}, nil
		},
		data.TypeString: func(rawBool map[string]interface{}) (interface{}, error) {
			if rawBoolValue, ok := rawBool["value"]; ok {
				if s, ok := rawBoolValue.(string); ok {
					return data.HollowString{Value: s}, nil
				} else {
					return nil, errors.New("`value` field of string is not of type string")
				}
			}
			return data.HollowBool{}, nil
		},
		data.TypeUint: func(_ map[string]interface{}) (interface{}, error) {
			return nil, errors.New("type `uint` is only used for simple type syntax")
		},
		data.TypeInt:     parseHollowInt,
		data.TypeFloat:   parseHollowFloat,
		data.TypeList:    parseHollowList,
		data.TypeChoice:  parseHollowChoice,
		data.TypeMap:     parseHollowMap,
		data.TypeStruct:  parseHollowMap,
		data.TypeTime:    parseHollowTime,
		data.TypeSize:    parseHollowSize,
		data.TypeChoiceN: parseHollowChoiceN,
	}
}
