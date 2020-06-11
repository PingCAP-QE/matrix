package data

import "time"

type Time struct {
	time time.Duration
}

func NewTime(t int64) Time {
	return Time{time: time.Duration(t)}
}

func NewTimeFromString(s string) (*Time, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return nil, err
	}
	return &Time{time: d}, nil
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.time.String()), nil
}

func (t Time) Nanoseconds() int64 { return t.time.Nanoseconds() }
