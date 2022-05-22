package services

// IServiceManager defines a basic way to interface with services.
type IServiceManager interface {
	Start(retry bool) error
	Stop(retry bool) error
	Restart(retry bool) error
}

type IStore[T any] interface {
	Get(key string) (*T, error)
	GetAll() (map[string]*T, error)
	Put(key string, payload *T) error
}
