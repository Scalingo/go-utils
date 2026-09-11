package stringparameter

type Service interface {
	Handle(string) error
}
