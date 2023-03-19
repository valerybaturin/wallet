// Package operation has a business logic for our application.
package operation

import (
	"context"
	"log"

	"wallet/app/queue"
)

// WalletService has
type WalletService struct {
	store    Store
	producer queue.Service
}

// NewWalletService ...
func NewWalletService(store Store, producer queue.Service) *WalletService {
	return &WalletService{
		store:    store,
		producer: producer,
	}
}

// Deposit amount from request to the wallets's balance.
func (s *WalletService) Deposit(ctx context.Context, id string, req Request) error {
	err := s.store.Deposit(ctx, id, req.Amount)
	if err != nil {
		return err
	}

	// make a channel to catch an error from the goroutine.
	errCh := make(chan error, 1)
	defer close(errCh)

	// publish message to the queue.
	go s.producer.Operation("Wallet_Deposited", req.Amount, errCh)

	// catch error from the channel.
	err = <-errCh
	if err != nil {
		log.Printf("queue.Publish error: %s", err.Error())
	}

	return nil
}

// Withdraw amount from wallet's balance.
func (s *WalletService) Withdraw(ctx context.Context, id string, req Request) error {
	err := s.store.Withdraw(ctx, id, req.Amount)
	if err != nil {
		return err
	}

	// make a channel to catch an error from the goroutine.
	errCh := make(chan error, 1)
	defer close(errCh)

	// publish message to the queue.
	go s.producer.Operation("Wallet_Withdrawn", req.Amount, errCh)

	// catch error from the channel.
	err = <-errCh
	if err != nil {
		log.Printf("queue.Publish error: %s", err.Error())
	}

	return nil
}

// Transfer money from one wallet to another.
func (s *WalletService) Transfer(ctx context.Context, id string, req TransferRequest) error {
	err := s.store.Transfer(ctx, id, req)
	if err != nil {
		return err
	}

	// make a channel to catch an error from the goroutine.
	errCh := make(chan error, 1)
	defer close(errCh)

	// publish message to the queue.
	go s.producer.Operation("Wallet_Transfered", req.Amount, errCh)

	// catch error from the channel.
	err = <-errCh
	if err != nil {
		log.Printf("queue.Publish error: %s", err.Error())
	}

	return nil
}
