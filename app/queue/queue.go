package queue

// Service has all methods for queue publisher.
type Service interface {
	Wallet(string, chan error)
	Operation(string, float64, chan error)
}
