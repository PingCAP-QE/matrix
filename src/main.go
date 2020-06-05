package main

import (
	"flag"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"chaos-mesh/matrix/pkg/context"
	"chaos-mesh/matrix/pkg/node"
	"chaos-mesh/matrix/pkg/parser"
)

var (
	conf string
)

func main() {
	flag.StringVar(&conf, "c", "", "config file")
	flag.Parse()

	if conf == "" {
		panic("config file not provided")
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
		panic("file not valid, detail error message is supposed to be shown above")
	}

	values := ctx.Gen()
	for config, concrete := range values.Configs {
		config.Serializer.Dump(concrete, config.Target)
	}
}
