package data

// ConcreteInterface value
type ConcreteInterface interface {
	ConcreteValue() interface{}
}

type ConcreteInt struct{ Value int }

func (c ConcreteInt) ConcreteValue() interface{} { return c.Value }

var _ ConcreteInterface = (*ConcreteInt)(nil)
