package stringparameter

type Service interface {
	Handle(value string) error
}
