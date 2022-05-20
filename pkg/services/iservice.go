package services

// IServiceManager defines a basic way to interface with services.
type IServiceManager interface {
	Start(retry bool) error
	Stop(retry bool) error
	Restart(retry bool) error
}
