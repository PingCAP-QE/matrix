package parser

import (
	"errors"
	"fmt"

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
		log.L().Warn(fmt.Sprintf("%s not handled, return HollowInt\n", rawValue))
		return data.NewHollowInt(utils.MinInt, utils.MaxInt), nil
	}
}

func parseMap(rawMap map[string]interface{}) (interface{}, error) {
	var hollowValue interface{}
	var err error
	if hollowType, ok := rawMap["type"]; ok {
		hollowValue, err = parseHollowValue(rawMap)
		if err == nil {
			return hollowValue, nil
		} else {
			log.L().Warn(fmt.Sprintf("HollowType %s parse failed, fall back to HollowMap", hollowType))
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
	hollowChoice.List = make([]data.HollowBranch, len(rawList))
	for i, k := range rawList {
		hollowBranch, err := parseBranch(k)
		if err != nil {
			return nil, err
		}
		hollowChoice.List[i] = *hollowBranch
	}
	return hollowChoice, nil
}

func parseBranch(rawBranch interface{}) (*data.HollowBranch, error) {
	return nil, errors.New("`parseBranch` not implemented")
}

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
	default:
		return data.HollowString{}, nil
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

func parseIntRange(raw interface{}) ([]int, error) {
	rangeList := raw.([]interface{})
	dur := make([]int, len(rangeList))
	for i, v := range rangeList {
		switch v.(type) {
		case float64:
			dur[i] = int(v.(float64))
		case string:
			return nil, errors.New(ExprNotSupportedMessage)
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
