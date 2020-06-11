package random

import (
	"chaos-mesh/matrix/pkg/node/data"
	"chaos-mesh/matrix/pkg/utils"
	"fmt"
	"math/rand"

	"github.com/c2h5oh/datasize"
)

func randUIntN(ui uint) uint {
	if ui <= uint(utils.MaxInt) {
		return uint(rand.Intn(int(ui)))
	} else {
		return uint(rand.Intn(utils.MaxInt)) + randUIntN(ui-uint(utils.MaxInt))
	}
}

func RandInt(start int, end int) int {
	if start == end {
		return start
	} else if start < end {
		return int(randUIntN(uint(end-start))) + start
	} else {
		panic("`RandInt` get a larger start than end")
	}
}

func RandFloat(start float64, end float64) float64 {
	if start == end {
		return start
	} else if start < end {
		return rand.Float64()*(start-end) + end
	} else {
		panic("`RandFloat` get a larger start than end")
	}
}

func RandChoose(brs []interface{}) interface{} {
	return brs[randUIntN(uint(len(brs)))]
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
		panic(fmt.Sprintf("time range too wide: %s - %s", start, end))
	}

}
