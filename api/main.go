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

package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"
	"github.com/pingcap/log"

	"github.com/chaos-mesh/matrix/pkg/context"
	"github.com/chaos-mesh/matrix/pkg/node"
	"github.com/chaos-mesh/matrix/pkg/parser"
)

func Gen(matrixConfig, output string, seed int64) error {
	cont, err := ioutil.ReadFile(matrixConfig)
	if err != nil {
		return err
	}

	var body node.MatrixConfigFile
	err = yaml.Unmarshal(cont, &body)
	if err != nil {
		return err
	}

	var ctx *context.MatrixContext
	ctx, err = parser.ParseFile(body)

	if err != nil {
		return errors.New(fmt.Sprintf("file not valid: %s", err.Error()))
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	log.L().Info(fmt.Sprintf("Matrix SEED: %d", seed))

	err = ctx.Dump(seed, output)
	if err != nil {
		return err
	}
	return nil
}
