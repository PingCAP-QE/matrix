package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pingcap/log"

	"chaos-mesh/matrix/pkg/context"
	"chaos-mesh/matrix/pkg/node"
	"chaos-mesh/matrix/pkg/parser"
)

var (
	conf   string
	output string
	err    error
)

func main() {
	flag.StringVar(&conf, "c", "", "config file")
	flag.StringVar(&output, "d", ".", "output folder")
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

	cont, err := ioutil.ReadFile(conf)
	if err != nil {
		panic(err)
	}

	var body node.MatrixConfigFile
	err = yaml.Unmarshal(cont, &body)
	if err != nil {
		panic(err)
	}

	var ctx *context.MatrixContext
	ctx, err = parser.ParseFile(body)

	if err != nil {
		panic(fmt.Sprintf("file not valid: %s", err.Error()))
	}

	values := ctx.Gen()
	for config, concrete := range values.Configs {
		err = config.Serializer.Dump(concrete, path.Join(output, config.Target))
		if err != nil {
			fmt.Printf("Error %s when dumping %s", err.Error(), concrete)
		}
	}
}
