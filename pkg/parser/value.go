package parser

import (
	"errors"
	"fmt"
	"github.com/pingcap/log"

	"chaos-mesh/matrix/pkg/node/data"
	"chaos-mesh/matrix/pkg/utils"
)

func parseTree(rawValue interface{}) (interface{}, error) {
	var err error

	switch rawValue.(type) {
	case map[string]interface{}:
		mapValue := rawValue.(map[string]interface{})
		if hollowType, ok := mapValue["type"]; ok {
			// full decl with type
			switch hollowType {
			case data.TypeInt:
				return parseInt(mapValue)
			}
		} else {
			// map
			var hollowMap data.HollowMap
			hollowMap.Map = make(map[string]interface{})
			for k := range mapValue {
				var v interface{}
				v, err = parseTree(mapValue[k])
				if err != nil {
					return nil, err
				}
				hollowMap.Map[k] = v
			}
			return hollowMap, nil
		}
	case []interface{}:
		// list
	case string:
		switch rawValue.(string) {
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
	case float64:
		// int / float literal
		floatValue := rawValue.(float64)
		intValue := int(floatValue)
		if floatValue == float64(intValue) {
			return data.NewHollowInt(intValue, intValue), nil
		} else {
			// todo: is a float
		}
	}
	log.L().Warn(fmt.Sprintf("%s not handled, return HollowInt\n", rawValue))
	return data.NewHollowInt(utils.MinInt, utils.MaxInt), nil
}

func parseInt(raw map[string]interface{}) (interface{}, error) {
	var hollowInt data.HollowInt
	var err error

	if intRange, ok := raw["range"]; ok {
		hollowInt.Range, err = parseIntRange(intRange)
		if err != nil {
			if err.Error() != ExprNotSupportedMessage {
				return nil, err
			}
		}
	}
	return hollowInt, nil
}

var ExprNotSupportedMessage = "expr not supported"

func parseIntRange(raw interface{}) ([]int, error) {
	rangeList := raw.([]interface{})
	dur := make([]int, len(rangeList))
	for i, v := range rangeList {
		switch v.(type) {
		case float64:
			dur[i] = int(v.(float64))
		case string:
			// todo: parse expression
			return nil, errors.New(ExprNotSupportedMessage)
		default:
			return nil, errors.New(fmt.Sprintf("%s cannot be parsed as int", v))
		}
	}
	return dur, nil
}
