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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chaos-mesh/matrix/api"

	"github.com/pingcap/log"
)

var (
	conf   string
	output string
	seed   int64
	err    error
)

func main() {
	flag.StringVar(&conf, "c", "", "config file")
	flag.StringVar(&output, "d", ".", "output folder")
	flag.Int64Var(&seed, "s", 0, "seed of rand, default UTC nanoseconds of now")
	flag.Parse()

	if conf == "" {
		panic("config file not provided")
	}
	output, err = filepath.Abs(output)
	if err != nil {
		panic(fmt.Sprintf("output folder invalid: %s", output))
	}
	if stat, err := os.Stat(output); os.IsNotExist(err) || !stat.IsDir() {
		panic(fmt.Sprintf("output folder does not exist: %s", output))
	} else {
		log.L().Info(fmt.Sprintf("dumpping to %s", output))
	}

	err := api.Gen(conf, output, seed)
	if err != nil {
		panic(err)
	}
}
