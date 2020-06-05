package serializer

type Config struct {
	Target     string
	Serializer Serializer
}

type Serializer interface {
	Dump(value interface{}, target string)
}
