package parser

import (
	"errors"

	"chaos-mesh/matrix/pkg/context"
	"chaos-mesh/matrix/pkg/node"
	"chaos-mesh/matrix/pkg/serializer"
)

func ParseFile(rawData node.MatrixConfigFile) (*context.MatrixContext, error) {
	var ctx context.MatrixContext

	for _, config := range rawData.Configs {
		conf, err := ParseConfig(config)
		if err != nil {
			return nil, err
		}
		if _, exist := ctx.Configs[conf.Tag]; exist {
			return nil, errors.New("%s appear twice in config")
		}
		ctx.Configs[conf.Tag] = *conf
	}
	return &ctx, nil
}

func ParseConfig(rawConfig node.RawConfig) (*node.AbstractConfig, error) {
	var conf node.AbstractConfig
	var err error

	conf.Tag = rawConfig.Tag

	if err = parseSerializerConfig(rawConfig, &conf.Config); err != nil {
		return nil, err
	}
	if err = parseValue(rawConfig.Value, &conf.Hollow); err != nil {
		return nil, err
	}

	return &conf, nil
}

func parseSerializerConfig(_rawConfig node.RawConfig, _out *serializer.Config) error {
	return errors.New("not implemented")
}
