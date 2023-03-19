// Package memory contains all implementation to work
// with wallet data store.
package memory

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"wallet/app/oops"
	"wallet/app/storage"
	"wallet/app/wallet"
)

// Storage contains map to store data and
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

// Wallets returns all data from the storage.
func (s *Storage) Wallets(ctx context.Context) ([]wallet.Wallet, error) {
	s.RLock()
	defer s.RUnlock()

	if len(s.data) == 0 {
		return nil, nil
	}

	wallets := make([]wallet.Wallet, 0, len(s.data))
	for k, v := range s.data {
		wallets = append(wallets, wallet.Wallet{
			ID:      k,
			Name:    v.Name,
			Status:  v.Status,
			Balance: v.Balance,
		})
	}

	return wallets, nil
}

// Wallet finds one record in map and returns it to the service.
func (s *Storage) Wallet(ctx context.Context, id string) (wallet.Wallet, error) {
	s.RLock()
	defer s.RUnlock()

	wal, found := s.data[id]
	if !found {
		return wallet.Wallet{}, oops.ErrNotFound
	}

	return wallet.Wallet{
		ID:      id,
		Name:    wal.Name,
		Balance: wal.Balance,
		Status:  wal.Status,
	}, nil
}

// CreateWallet generates id and stores new wallet into the map.
func (s *Storage) CreateWallet(ctx context.Context, req wallet.Request) (wallet.Wallet, error) {
	// generate id with length 8
	id := generateID(8)

	s.Lock()
	defer s.Unlock()

	s.data[id] = storage.Wallet{
		Name:   req.Name,
		Status: "active",
	}

	return wallet.Wallet{
		ID:     id,
		Name:   req.Name,
		Status: "active",
	}, nil
}

// UpdateWallet updates name in the map.
func (s *Storage) UpdateWallet(ctx context.Context, req wallet.Request, id string) error {
	s.Lock()
	defer s.Unlock()

	wal, found := s.data[id]
	if !found || wal.Status == "inactive" {
		return oops.ErrNotFound
	}

	s.data[id] = storage.Wallet{
		Name: req.Name,
	}

	return nil
}

// DeleteWallet removes wallet from the map.
func (s *Storage) DeleteWallet(ctx context.Context, id string) error {
	s.RLock()
	defer s.RUnlock()

	wal, found := s.data[id]
	if !found || wal.Status == "inactive" {
		return oops.ErrNotFound
	}

	s.data[id] = storage.Wallet{
		Status: "inactive",
	}

	return nil
}

func generateID(n int) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		// TODO: change to 1.20
		rand.Seed(time.Now().UnixNano())
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
