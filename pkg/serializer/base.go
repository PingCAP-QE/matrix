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

package serializer

import (
	"errors"
	"fmt"
)

type Config struct {
	Target     string
	Serializer Serializer
}

type Serializer interface {
	Dump(value interface{}, target string) error
}

func ParseSerializerName(name string) (Serializer, error) {
	serializer, ok := serializerMap[name]
	if ok {
		return serializer, nil
	}
	return nil, errors.New(fmt.Sprintf("%s not a valid serializer name", name))
}

var serializerMap = map[string]Serializer{
	"toml": TomlSerializer{},
	"yaml": YamlSerializer{},
	"stmt": StatementSerializer{},
}
