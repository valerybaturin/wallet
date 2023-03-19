package operation

import (
	"encoding/json"
	"net/http"

	"wallet/app/oops"
	"wallet/app/response"

	"github.com/go-chi/chi/v5"
)

// Handler contains app Service and a router.
type Handler struct {
	router    *chi.Mux
	operation WalletService
}

// NewHandler is a constructor which accepts operation Service and
// returns a pointer to the Handler.
func NewHandler(router *chi.Mux, service WalletService) *Handler {
	return &Handler{
		router:    router,
		operation: service,
	}
}

// Register operation routes.
func (h *Handler) Register() {
	h.router.Group(func(r chi.Router) {
		r.Post("/wallets/{id}/deposit", h.deposit)
		r.Post("/wallets/{id}/withdraw", h.withdraw)
		r.Post("/wallets/{id}/transfer", h.transfer)
	})
}

func (h *Handler) deposit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.OperationError(w, http.StatusBadRequest, oops.ErrBadReqMessage, 0)
		return
	}

	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		response.OperationError(w, http.StatusBadRequest, oops.ErrBadReqMessage, requestBody.Amount)
		return
	}

	err = h.operation.Deposit(r.Context(), id, requestBody)
	if err != nil {
		response.OperationError(w, http.StatusInternalServerError, oops.ErrIntServMessage, requestBody.Amount)
		return
	}

	response.OperationSuccess(w, http.StatusOK, requestBody.Amount)
}

func (h *Handler) withdraw(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.OperationError(w, http.StatusBadRequest, oops.ErrBadReqMessage, 0)
		return
	}

	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		response.OperationError(w, http.StatusBadRequest, oops.ErrBadReqMessage, requestBody.Amount)
		return
	}

	err = h.operation.Withdraw(r.Context(), id, requestBody)
	if err != nil {
		response.OperationError(w, http.StatusInternalServerError, oops.ErrIntServMessage, requestBody.Amount)
		return
	}

	response.OperationSuccess(w, http.StatusOK, requestBody.Amount)
}

func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.OperationError(w, http.StatusBadRequest, oops.ErrBadReqMessage, 0)
		return
	}

	var requestBody TransferRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		response.OperationError(w, http.StatusBadRequest, oops.ErrBadReqMessage, requestBody.Amount)
		return
	}

	err = h.operation.Transfer(r.Context(), id, requestBody)
	if err != nil {
		response.OperationError(w, http.StatusInternalServerError, oops.ErrIntServMessage, requestBody.Amount)
		return
	}

	response.OperationSuccess(w, http.StatusOK, requestBody.Amount)
}
