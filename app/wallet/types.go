package wallet

import "context"

// Wallet contains all fields to define wallet.
type Wallet struct {
	ID      string  `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Balance float64 `json:"balance,omitempty"`
	Status  string  `json:"status,omitempty"`
}

// Request contains fields for client request.
type Request struct {
	Name string `json:"name" validate:"required,gte=1"`
}

// Store contains all methods to store data into the storage.
type Store interface {
	Wallets(context.Context) ([]Wallet, error)
	Wallet(context.Context, string) (Wallet, error)
	CreateWallet(context.Context, Request) (Wallet, error)
	UpdateWallet(context.Context, Request, string) error
	DeleteWallet(context.Context, string) error
}

// Service contains all methods from wallet service.
type Service interface {
	List(context.Context) ([]Wallet, error)
	Item(context.Context, string) (Wallet, error)
	Create(context.Context, Request) (Wallet, error)
	Update(context.Context, Request, string) error
	Delete(context.Context, string) error
}
