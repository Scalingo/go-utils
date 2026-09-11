package service

type Embedded interface {
	EmbeddedMethod()
}

type Service interface {
	Embedded
	Ping()
	Transform(name string, payload []byte) (int, error)
}
