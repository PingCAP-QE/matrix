package parser

import (
	"errors"

	"chaos-mesh/matrix/pkg/context"
	"chaos-mesh/matrix/pkg/node"
	"chaos-mesh/matrix/pkg/serializer"
)

func ParseFile(rawData node.MatrixConfigFile) (*context.MatrixContext, error) {
	var ctx context.MatrixContext

	ctx.Configs = make(map[string]node.AbstractConfig)

	for _, config := range rawData.Configs {
		conf, err := parseConfig(config)
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

func parseConfig(rawConfig node.RawConfig) (*node.AbstractConfig, error) {
	var conf node.AbstractConfig
	var err error

	conf.Tag = rawConfig.Tag

	if err = parseSerializerConfig(rawConfig, &conf.Config); err != nil {
		return nil, err
	}
	if conf.Hollow, err = parseTree(rawConfig.Value); err != nil {
		return nil, err
	}

	return &conf, nil
}

func parseSerializerConfig(rawConfig node.RawConfig, out *serializer.Config) error {
	var err error
	out.Serializer, err = serializer.ParseSerializerName(rawConfig.Serializer)
	if err != nil {
		return err
	}
	out.Target = rawConfig.Target
	return nil
}
