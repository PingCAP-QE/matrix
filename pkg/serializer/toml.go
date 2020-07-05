package serializer

import (
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

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
