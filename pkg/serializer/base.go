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
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
)

type Config struct {
	Target     string
	Serializer Serializer
}

type Serializer interface {
	Dump(value interface{}, target string) error
}

var serializerMap = map[string]Serializer{
	"toml": TomlSerializer{},
	"yaml": YamlSerializer{},
}

func ParseSerializerName(name string) (Serializer, error) {
	serializer, ok := serializerMap[name]
	if ok {
		return serializer, nil
	}
	return nil, errors.New(fmt.Sprintf("%s not a valid serializer name", name))
}

type TomlSerializer struct{}
type YamlSerializer struct{}

func (s TomlSerializer) Dump(value interface{}, target string) error {
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(f)
	encoder.Indent = ""
	return encoder.Encode(value)
}

func (s YamlSerializer) Dump(value interface{}, target string) error {
	var text []byte
	var err error
	var f *os.File
	text, err = yaml.Marshal(value)
	if err != nil {
		return err
	}
	f, err = os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(text)
	return err
}
