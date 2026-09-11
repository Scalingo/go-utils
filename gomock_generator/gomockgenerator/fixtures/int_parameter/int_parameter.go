package intparameter

type Service interface {
	Handle(value int) error
}
