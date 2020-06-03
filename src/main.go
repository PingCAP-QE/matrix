package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
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

	var body interface{}
	err = yaml.Unmarshal(cont, &body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", body)
}
