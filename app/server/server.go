// Package server initializing and runs a web server on static port 3000
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"wallet/app/operation"
	operStorage "wallet/app/operation/memory"
	"wallet/app/queue"
	"wallet/app/storage"
	"wallet/app/wallet"
	walletStorage "wallet/app/wallet/memory"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

// Server contains http.Server.
type Server struct {
	Router *chi.Mux
	Queue  *queue.NSQ
	HTTP   *http.Server
}

// New is a constructor which initializes new Server.
func New(queue *queue.NSQ) *Server {
	r := chi.NewRouter()

	return &Server{
		Router: r,
		Queue:  queue,
		HTTP: &http.Server{
			Addr:         ":3000",           // app port
			Handler:      r,                 // set the default handler
			ReadTimeout:  5 * time.Second,   // max time to read request from the client
			WriteTimeout: 10 * time.Second,  // max time to write response to the client
			IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		},
	}
}

// SetupMiddleware register middlewares.
func (s *Server) SetupMiddleware() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
}

// SetupApp registers app services.
func (s *Server) SetupApp() {
	storage := storage.NewMemory()
	walletStore := walletStorage.NewStorage(storage.Data)

	walletService := wallet.NewAppService(walletStore, s.Queue)
	walletHandler := wallet.NewHandler(s.Router, *walletService)
	walletHandler.Register()

	operStore := operStorage.NewStorage(storage.Data)
	operationService := operation.NewWalletService(operStore, s.Queue)
	operationHandler := operation.NewHandler(s.Router, *operationService)
	operationHandler.Register()
}

// Start runs HTTP server.
func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errs, ctx := errgroup.WithContext(ctx)

	log.Println("starting web server on port 3000")

	errs.Go(func() error {
		if err := s.HTTP.ListenAndServe(); err != nil {
			return fmt.Errorf("listen and serve error: %w", err)
		}
		return nil
	})

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully")

	// Perform application shutdown with a maximum timeout of 5 seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.HTTP.Shutdown(timeoutCtx); err != nil {
		log.Println(err.Error())
	}

	return nil
}
