package parser

import (
	"errors"
	"fmt"
	"time"

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
		return data.NewHollowInt(utils.MinInt, utils.MaxInt), nil
	}
}

func parseMap(rawMap map[string]interface{}) (interface{}, error) {
	var hollowValue interface{}
	var err error
	if hollowType, ok := rawMap["type"]; ok {
		if _, ok := rawMap["when"]; ok {
			// todo: handle when condition
			_, _ = parseCondition(rawMap)
			delete(rawMap, "when")
		}
		hollowValue, err = parseHollowValue(rawMap)
		if err == nil {
			return hollowValue, nil
		} else {
			log.L().Warn(fmt.Sprintf("type %s parse failed with message \"%s\", fall back to HollowMap", hollowType, err.Error()))
		}
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

func parseCondition(rawMap map[string]interface{}) (*data.HollowCondition, error) {
	if cond, ok := rawMap["when"]; ok {
		log.L().Warn(fmt.Sprintf("ignoring constraint: %s", cond))
	}
	return nil, nil
}

// parseString convert string literals to
func parseString(rawString string) (interface{}, error) {
	switch rawString {
	case data.TypeBool:
		return data.HollowBool{}, nil
	case data.TypeInt:
		return data.NewHollowInt(utils.MinInt, utils.MaxInt), nil
	case data.TypeUint:
		return data.NewHollowInt(0, utils.MaxInt), nil
	case data.TypeString:
		return data.HollowString{}, nil
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
		return data.NewHollowInt(intValue, intValue)
	} else {
		return data.NewHollowFloat(rawFloat, rawFloat)
	}
}

func parseHollowValue(rawHollow map[string]interface{}) (interface{}, error) {
	switch rawHollow["type"] {
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
		hollowInt.Range, err = parseIntRange(intRange)
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
		hollowFloat.Range, err = parseFloatRange(floatRange)
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
		hollowTime.Range, err = parseTimeRange(timeRange)
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
		hollowSize.Range, err = parseSizeRange(timeRange)
		if err != nil {
			return nil, err
		}
	}
	return hollowSize, nil
}

func parseIntRange(raw interface{}) ([]int, error) {
	rangeList := raw.([]interface{})
	dur := make([]int, len(rangeList))
	for i, v := range rangeList {
		switch v.(type) {
		case float64:
			dur[i] = int(v.(float64))
		case string:
			return nil, errors.New(fmt.Sprintf("%s: %s", ExprNotSupportedMessage, v.(string)))
		default:
			return nil, errors.New(fmt.Sprintf("%s cannot be parsed as int", v))
		}
	}
	return dur, nil
}

func parseFloatRange(raw interface{}) ([]float64, error) {
	rangeList := raw.([]interface{})
	dur := make([]float64, len(rangeList))
	for i, v := range rangeList {
		switch v.(type) {
		case float64:
			dur[i] = v.(float64)
		case string:
			return nil, errors.New(ExprNotSupportedMessage)
		default:
			return nil, errors.New(fmt.Sprintf("%s cannot be parsed as float", v))
		}
	}
	return dur, nil
}

func parseTimeRange(raw interface{}) ([]time.Duration, error) {
	var err error
	rangeList := raw.([]interface{})
	dur := make([]time.Duration, len(rangeList))
	for i, v := range rangeList {
		if timeStr, ok := v.(string); ok {
			dur[i], err = time.ParseDuration(timeStr)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(fmt.Sprintf("%s cannot be parsed as time", v))
		}
	}
	return dur, nil
}

func parseSizeRange(raw interface{}) ([]datasize.ByteSize, error) {
	var err error
	rangeList := raw.([]interface{})
	dur := make([]datasize.ByteSize, len(rangeList))
	for i, v := range rangeList {
		if sizeStr, ok := v.(string); ok {
			var size datasize.ByteSize
			err = size.UnmarshalText([]byte(sizeStr))
			if err != nil {
				return nil, err
			}
			dur[i] = size
		} else {
			return nil, errors.New(fmt.Sprintf("%s cannot be parsed as size", v))
		}
	}
	return dur, nil
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
