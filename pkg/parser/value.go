package parser

import (
	"chaos-mesh/matrix/pkg/node/data"
	"errors"
)

func parseValue(rawValue interface{}) (interface{}, error) {
	var err error

	switch rawValue.(type) {
	case map[string]interface{}:
		mapValue := rawValue.(map[string]interface{})
		if hollowType, ok := mapValue["type"]; ok {
			switch hollowType {
			case data.TypeInt:
				return parseInt(mapValue)
			}
		} else {
			var hollowMap data.HollowMap
			for k := range mapValue {
				hollowMap.Map = make(map[string]interface{})
				hollowMap.Map[k], err = parseValue(mapValue[k])
				if err != nil {
					return nil, err
				}
			}
			return hollowMap, nil
		}
	}
	return data.HollowInt{}, nil
}

func parseInt(raw map[string]interface{}) (interface{}, error) {
	var hollowInt data.HollowInt
	if intRange, ok := raw["range"]; ok {
		_, err := parseRange(intRange)
		return nil, err
	}
	return hollowInt, nil
}

func parseRange(raw interface{}) (interface{}, error) {
	return nil, errors.New("parseRange not implemented")
}
