package random

import (
	"chaos-mesh/matrix/pkg/utils"
	"math/rand"
	"time"

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
	case size > datasize.EB:
		return datasize.EB
	case size > datasize.PB:
		return datasize.PB
	case size > datasize.TB:
		return datasize.TB
	case size > datasize.GB:
		return datasize.GB
	case size > datasize.MB:
		return datasize.MB
	case size > datasize.KB:
		return datasize.KB
	default:
		return 1
	}
}

func RandSize(start datasize.ByteSize, end datasize.ByteSize) datasize.ByteSize {
	startUnit := uint64(sizeUnit(start))

	startInt, endInt := int(uint64(start)/startUnit), int(uint64(end)/startUnit)

	return datasize.ByteSize(uint64(RandInt(startInt, endInt)) * startUnit)
}

func timeUnit(t time.Duration) time.Duration {
	switch {
	case t > time.Hour:
		return time.Hour
	case t > time.Minute:
		return time.Minute
	case t > time.Second:
		return time.Second
	default:
		return time.Millisecond
	}
}

func RandTime(start time.Duration, end time.Duration) time.Duration {
	startUnit := int64(timeUnit(start))

	startInt, endInt := int(int64(start)/startUnit), int(int64(end)/startUnit)

	return time.Duration(int64(RandInt(startInt, endInt)) * startUnit)
}
