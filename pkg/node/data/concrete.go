package data

// ConcreteInterface value
type ConcreteInterface interface {
	ConcreteValue() interface{}
	ConcreteType() string
}

type ConcreteInt struct {
}
