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

package parser

import (
	"errors"

	"github.com/chaos-mesh/matrix/pkg/context"
	"github.com/chaos-mesh/matrix/pkg/node"
	"github.com/chaos-mesh/matrix/pkg/serializer"
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
