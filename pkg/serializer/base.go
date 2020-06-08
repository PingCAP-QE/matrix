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

var serializerMap = map[string]Serializer{
	"toml": TomlSerializer{},
}

func ParseSerializerName(name string) (Serializer, error) {
	serializer, ok := serializerMap[name]
	if ok {
		return serializer, nil
	}
	return nil, errors.New(fmt.Sprintf("%s not a valid serializer name", name))
}

type TomlSerializer struct{}

func (s TomlSerializer) Dump(value interface{}, target string) error {
	return nil
}
