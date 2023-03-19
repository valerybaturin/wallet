// Package memory contains all implementation to work
// with wallet operations in data store.
package memory

import (
	"context"
	"fmt"
	"sync"

	"wallet/app/oops"
	"wallet/app/operation"
	"wallet/app/storage"
)

// Storage contains map to work with data and
// RWMutex to sync read/write operations.
type Storage struct {
	data map[string]storage.Wallet
	sync.RWMutex
}

// NewStorage is a constructor for storage.
func NewStorage(data map[string]storage.Wallet) *Storage {
	return &Storage{
		data: data,
	}
}

// Deposit adds amount to the balance.
func (s *Storage) Deposit(ctx context.Context, id string, amount float64) error {
	s.Lock()
	defer s.Unlock()

	wallet, found := s.data[id]
	if !found || wallet.Status == "inactive" {
		return oops.ErrNotFound
	}

	s.data[id] = storage.Wallet{
		Balance: wallet.Balance + amount,
	}

	return nil
}

// Withdraw ...
func (s *Storage) Withdraw(ctx context.Context, id string, amount float64) error {
	s.Lock()
	defer s.Unlock()

	wallet, found := s.data[id]
	if !found || wallet.Status == "inactive" {
		return oops.ErrNotFound
	}

	if wallet.Balance == 0 || wallet.Balance < amount {
		return oops.ErrNotEnoMon
	}

	s.data[id] = storage.Wallet{
		Balance: wallet.Balance - amount,
	}

	return nil
}

// Transfer ...
func (s *Storage) Transfer(ctx context.Context, id string, data operation.TransferRequest) error {
	err := s.Withdraw(ctx, id, data.Amount)
	if err != nil {
		return fmt.Errorf("Transfer Withdraw error: %w", err)
	}

	err = s.Deposit(ctx, data.TransferTo, data.Amount)
	if err != nil {
		return fmt.Errorf("Transfer Deposit error: %w", err)
	}

	return nil
}
