package node

// raw config file

type MatrixConfigFile struct {
	Configs []RawConfig `yaml:"configs"`
}

type RawConfig struct {
	Tag        string      `yaml:"tag"`
	Serializer string      `yaml:"serializer"`
	Target     string      `yaml:"target"`
	Value      interface{} `yaml:"value"`
}
