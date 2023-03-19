package operation

import "context"

// Request contains fields for client request.
type Request struct {
	Amount float64 `json:"amount"`
}

// TransferRequest contains fields for client request.
type TransferRequest struct {
	Amount     float64 `json:"amount"`
	TransferTo string  `json:"transfer_to"`
}

// Store contains all methods to store data into the storage.
type Store interface {
	Deposit(context.Context, string, float64) error
	Withdraw(context.Context, string, float64) error
	Transfer(context.Context, string, TransferRequest) error
}

// Service contains all methods from operation service.
type Service interface {
	Deposit(context.Context, string, Request) error
	Withdraw(context.Context, string, Request) error
	Transfer(context.Context, string, TransferRequest) error
}
