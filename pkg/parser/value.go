package parser

import (
	"errors"
	"fmt"

	"github.com/c2h5oh/datasize"

	"github.com/pingcap/log"

	"chaos-mesh/matrix/pkg/node/data"
	"chaos-mesh/matrix/pkg/utils"
)

const ExprNotSupportedMessage = "expr not supported"

func parseTree(rawValue interface{}) (interface{}, error) {
	switch rawValue.(type) {
	case map[string]interface{}:
		return parseMap(rawValue.(map[string]interface{}))
	case []interface{}:
		return parseChoice(rawValue.([]interface{}))
	case string:
		return parseString(rawValue.(string))
	case float64:
		return parseFloat(rawValue.(float64)), nil
	default:
		log.L().Warn(fmt.Sprintf("%s not handled, return HollowInt", rawValue))
		return data.HollowInt{RangeStart: utils.MinInt, RangeEnd: utils.MaxInt}, nil
	}
}

func parseMap(rawMap map[string]interface{}) (interface{}, error) {
	var err error
	if _, ok := rawMap["type"]; ok {
		var hollowValue interface{}
		var rawCondition interface{}

		if rawCondition, ok = rawMap["when"]; ok {
			delete(rawMap, "when")
		}

		hollowValue, err = parseHollowValue(rawMap)
		if err != nil {
			return nil, err
		}

		if rawCondition != nil {
			var hollowCondition *data.HollowCondition
			hollowCondition, err = parseCondition(rawCondition)
			if err != nil {
				return nil, err
			}
			hollowValue.(data.HollowInterface).SetCondition(hollowCondition)
		}

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

// parseString convert string literals to corresponding values, default value will be handled here
func parseString(rawString string) (interface{}, error) {
	switch rawString {
	case data.TypeBool:
		return data.HollowBool{}, nil
	case data.TypeInt:
		return data.HollowInt{RangeStart: utils.MinInt, RangeEnd: utils.MaxInt}, nil
	case data.TypeUint:
		return data.HollowInt{RangeStart: 0, RangeEnd: utils.MaxInt}, nil
	case data.TypeFloat:
		return data.HollowFloat{RangeStart: 0, RangeEnd: 1}, nil
	case data.TypeString:
		return data.HollowString{Value: ""}, nil
	case data.TypeTime:
		return data.HollowTime{}, nil
	case data.TypeSize:
		return data.HollowSize{}, nil
	default:
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

func parseHollowValue(rawHollow map[string]interface{}) (interface{}, error) {
	switch rawHollow["type"] {
	case data.TypeBool:
		return data.HollowBool{}, nil
	case data.TypeUint:
		return nil, errors.New("type `uint` is only used for simple type syntax")
	case data.TypeInt:
		return parseHollowInt(rawHollow)
	case data.TypeFloat:
		return parseHollowFloat(rawHollow)
	case data.TypeList:
		return parseHollowList(rawHollow)
	case data.TypeMap, data.TypeStruct:
		return parseHollowMap(rawHollow)
	case data.TypeTime:
		return parseHollowTime(rawHollow)
	case data.TypeSize:
		return parseHollowSize(rawHollow)
	}
	return nil, errors.New(fmt.Sprintf("parseHollowValue for type %s not implemented", rawHollow["type"]))
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
