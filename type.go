package gim

type withStatus interface {
	Status() int
}
type withBody interface {
	Body() interface{}
}
