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
