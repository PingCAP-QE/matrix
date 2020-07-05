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

package random

import (
	"fmt"
	"math/rand"

	"github.com/c2h5oh/datasize"

	"chaos-mesh/matrix/pkg/node/data"
	"chaos-mesh/matrix/pkg/utils"
)

const selectRangeBoundaryProb = 0.05

func Seed(seed int64) {
	rand.Seed(seed)
}

func randUIntN(ui uint) uint {
	if ui <= uint(utils.MaxInt) {
		return uint(rand.Intn(int(ui)))
	} else {
		return uint(rand.Intn(utils.MaxInt)) + randUIntN(ui-uint(utils.MaxInt))
	}
}

// return true with a probability of `prob`
func randProb(prob float64) bool {
	if prob < 0 {
		return false
	} else if prob >= 1 {
		return true
	} else {
		return prob > rand.Float64()
	}
}

func randAny(values ...interface{}) interface{} {
	return values[randUIntN(uint(len(values)))]
}

func RandInt(start int, end int) int {
	if start == end {
		return start
	} else if randProb(selectRangeBoundaryProb) {
		return randAny(start, end).(int)
	} else if start < end {
		return int(randUIntN(uint(end-start))) + start
	} else {
		panic("`RandInt` get a larger start than end")
	}
}

func RandFloat(start float64, end float64) float64 {
	if start == end {
		return start
	} else if randProb(selectRangeBoundaryProb) {
		return randAny(start, end).(float64)
	} else if start < end {
		return rand.Float64()*(start-end) + end
	} else {
		panic("`RandFloat` get a larger start than end")
	}
}

func RandChoose(brs []interface{}) interface{} {
	return brs[randUIntN(uint(len(brs)))]
}

func RandChooseN(brs []interface{}, n int) []interface{} {
	if n > len(brs) {
		panic(fmt.Sprintf("`RandChooseN` n out of bound: %d, %v", n, brs))
	}
	var idx, selectedIdx []int
	for i := 0; i < len(brs); i++ {
		idx = append(idx, i)
	}
	rand.Shuffle(len(idx), func(i, j int) { idx[i], idx[j] = idx[j], idx[i] })
	for i := 0; i < n; i++ {
		selectedIdx = append(selectedIdx, idx[i])
	}
	results := make([]interface{}, n)
	for i, idx := range selectedIdx {
		results[i] = brs[idx]
	}
	return results
}

func sizeUnit(size datasize.ByteSize) datasize.ByteSize {
	switch {
	case size >= datasize.EB:
		return datasize.EB
	case size >= datasize.PB:
		return datasize.PB
	case size >= datasize.TB:
		return datasize.TB
	case size >= datasize.GB:
		return datasize.GB
	case size >= datasize.MB:
		return datasize.MB
	case size >= datasize.KB:
		return datasize.KB
	default:
		return 1
	}
}

func RandSize(start datasize.ByteSize, end datasize.ByteSize) datasize.ByteSize {
	if randProb(selectRangeBoundaryProb) {
		return randAny(start, end).(datasize.ByteSize)
	}

	var unit uint64
	if start == 0 {
		unit = uint64(sizeUnit(end))
	} else {
		unit = uint64(sizeUnit(start))
	}

	startInt, endInt := int(uint64(start)/unit), int(uint64(end)/unit)

	return datasize.ByteSize(uint64(RandInt(startInt, endInt)) * unit)
}

func RandTime(start data.Time, end data.Time) data.Time {
	if randProb(selectRangeBoundaryProb) {
		return randAny(start, end).(data.Time)
	}

	startInt, endInt := start.Nanoseconds(), end.Nanoseconds()

	var unit int64 = 1
	for startInt > 10 {
		startInt /= 10
		endInt /= 10
		unit *= 10
	}

	if startInt <= int64(utils.MaxInt) && endInt <= int64(utils.MaxInt) {
		return data.NewTime(int64(RandInt(int(startInt), int(endInt))) * unit)
	} else {
		panic(fmt.Sprintf("time range too wide: %v - %v", start, end))
	}

}
