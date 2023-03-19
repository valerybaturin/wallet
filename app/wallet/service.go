// Package wallet has a business logic for our application.
package wallet

import (
	"context"
	"log"

	"wallet/app/queue"
)

// AppService contains Store interface.
type AppService struct {
	store    Store
	producer queue.Service
}

// NewAppService is a Service constructor.
func NewAppService(store Store, producer queue.Service) *AppService {
	return &AppService{
		store:    store,
		producer: producer,
	}
}

// List returns list of Wallets from the Store.
func (s *AppService) List(ctx context.Context) ([]Wallet, error) {
	// get wallets from the store.
	wallets, err := s.store.Wallets(ctx)
	if err != nil {
		return nil, err
	}

	if len(wallets) == 0 {
		return nil, nil
	}

	return wallets, nil
}

// Item returns a wallet from the store to the client.
func (s *AppService) Item(ctx context.Context, id string) (Wallet, error) {
	wallet, err := s.store.Wallet(ctx, id)
	if err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

// Create saves a new wallet into the storage.
func (s *AppService) Create(ctx context.Context, req Request) (Wallet, error) {
	wallet, err := s.store.CreateWallet(ctx, req)
	if err != nil {
		return Wallet{}, err
	}

	// make a channel to catch an error from the goroutine.
	errCh := make(chan error, 1)
	defer close(errCh)

	// publish message to the queue.
	go s.producer.Wallet("Wallet_Created", errCh)

	// catch error from the channel.
	err = <-errCh
	if err != nil {
		log.Printf("queue.Publish error: %s", err.Error())
	}

	return Wallet{
		ID:     wallet.ID,
		Name:   wallet.Name,
		Status: wallet.Status,
	}, nil
}

// Update updates the name in the storage.
func (s *AppService) Update(ctx context.Context, req Request, id string) error {
	err := s.store.UpdateWallet(ctx, req, id)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes wallet from the storage.
func (s *AppService) Delete(ctx context.Context, id string) error {
	err := s.store.DeleteWallet(ctx, id)
	if err != nil {
		return err
	}

	// make a channel to catch an error from the goroutine.
	errCh := make(chan error, 1)
	defer close(errCh)

	// publish message to the queue.
	go s.producer.Wallet("Wallet_Deleted", errCh)

	// catch error from the channel.
	err = <-errCh
	if err != nil {
		log.Printf("queue.Publish error: %s", err)
	}

	return nil
}
