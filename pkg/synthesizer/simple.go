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

package synthesizer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pingcap/log"

	"github.com/chaos-mesh/matrix/pkg/node/data"
	"github.com/chaos-mesh/matrix/pkg/random"
)

func SimpleRecGen(hollow interface{}) interface{} {
	switch hollow.(type) {
	case data.HollowBool:
		b := hollow.(data.HollowBool).Value
		if b != nil {
			var boolValue bool = *b
			return boolValue
		}
		return random.RandChoose([]interface{}{true, false})
	case data.HollowInt:
		return random.RandInt(hollow.(data.HollowInt).RangeStart, hollow.(data.HollowInt).RangeEnd)
	case data.HollowFloat:
		return random.RandFloat(hollow.(data.HollowFloat).RangeStart, hollow.(data.HollowFloat).RangeEnd)
	case data.HollowTime:
		return random.RandTime(hollow.(data.HollowTime).RangeStart, hollow.(data.HollowTime).RangeEnd)
	case data.HollowSize:
		return random.RandSize(hollow.(data.HollowSize).RangeStart, hollow.(data.HollowSize).RangeEnd)
	case data.HollowString:
		hollowString := hollow.(data.HollowString)
		if hollowString.Value != "" {
			return hollowString.Value
		} else {
			return "rand-string"
		}
	case data.HollowMap:
		res := make(map[string]interface{})
		for _, k := range sortedMapKey(hollow.(data.HollowMap).Map) {

			res[k] = SimpleRecGen(hollow.(data.HollowMap).Map[k])
		}
		return res
	case data.HollowList:
		var res []interface{}
		for _, v := range hollow.(data.HollowList).List {
			res = append(res, SimpleRecGen(v))
		}
		return res
	case data.HollowChoice:
		return SimpleRecGen(random.RandChoose(hollow.(data.HollowChoice).List))
	case data.HollowChoiceN:
		choiceN := hollow.(data.HollowChoiceN)
		n := choiceN.N
		if n == 0 {
			n = random.RandInt(1, len(choiceN.List))
		}
		results := make([]interface{}, n)
		strResults := make([]string, n)
		tryJoin := true
		for i, hollow := range random.RandChooseN(choiceN.List, n) {
			v := SimpleRecGen(hollow)
			results[i] = v
			if tryJoin {
				s, isStr := v.(string)
				tryJoin = tryJoin && isStr
				strResults[i] = s
			}
		}
		if tryJoin {
			return strings.Join(strResults, choiceN.Sep)
		} else {
			return results
		}
	default:
		log.L().Warn(fmt.Sprintf("unhandled value: %v", hollow))
		return fmt.Sprintf("%s", hollow)
	}
}

func sortedMapKey(rawMap map[string]interface{}) []string {
	var keys []string
	for s, _ := range rawMap {
		keys = append(keys, s)
	}
	sort.Strings(keys)
	return keys
}
