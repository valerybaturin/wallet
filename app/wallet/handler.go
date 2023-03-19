package wallet

import (
	"encoding/json"
	"errors"
	"net/http"

	"wallet/app/oops"
	"wallet/app/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler contains app Service and a router.
type Handler struct {
	router *chi.Mux
	wallet AppService
}

// NewHandler is a constructor which accepts app Service and
// returns a pointer to the Handler.
func NewHandler(router *chi.Mux, service AppService) *Handler {
	return &Handler{
		router: router,
		wallet: service,
	}
}

// Register wallet routes.
func (h *Handler) Register() {
	h.router.Group(func(r chi.Router) {
		r.Get("/wallets", h.list)
		r.Get("/wallets/{id}", h.item)
		r.Post("/wallet", h.create)
		r.Put("/wallets/{id}", h.update)
		r.Delete("/wallet/{id}", h.delete)
	})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	data, err := h.wallet.List(r.Context())
	if err != nil {
		response.WalletError(w, http.StatusInternalServerError, oops.ErrIntServMessage, "")
		return
	}
	if len(data) == 0 {
		response.WalletError(w, http.StatusNoContent, oops.ErrNoDataMessage, "")
		return
	}

	response.Data(w, http.StatusOK, data)
}

func (h *Handler) item(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WalletError(w, http.StatusBadRequest, oops.ErrBadReqMessage, "")
		return
	}

	data, err := h.wallet.Item(r.Context(), id)
	if err != nil {
		if errors.Is(err, oops.ErrNotFound) {
			response.WalletError(w, http.StatusNotFound, oops.ErrNotFoundMessage, id)
			return
		}
		response.WalletError(w, http.StatusInternalServerError, oops.ErrIntServMessage, id)
		return
	}

	response.Data(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		response.WalletError(w, http.StatusBadRequest, oops.ErrBadReqMessage, "")
		return
	}

	data, err := h.wallet.Create(r.Context(), requestBody)
	if err != nil {
		response.WalletError(w, http.StatusInternalServerError, oops.ErrIntServMessage, "")
		return
	}

	response.Data(w, http.StatusCreated, data)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WalletError(w, http.StatusBadRequest, oops.ErrBadReqMessage, id)
		return
	}

	var requestBody Request
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		response.WalletError(w, http.StatusBadRequest, oops.ErrBadReqMessage, id)
		return
	}
	if err = validator.New().Struct(&requestBody); err != nil {
		response.WalletError(w, http.StatusBadRequest, oops.ErrBadReqMessage, id)
		return
	}

	err = h.wallet.Update(r.Context(), requestBody, id)
	if err != nil {
		response.WalletError(w, http.StatusInternalServerError, oops.ErrIntServMessage, id)
		return
	}

	response.WalletSuccess(w, http.StatusOK, id)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WalletError(w, http.StatusBadRequest, oops.ErrBadReqMessage, id)
		return
	}

	err := h.wallet.Delete(r.Context(), id)
	if err != nil {
		response.WalletError(w, http.StatusInternalServerError, oops.ErrIntServMessage, id)
		return
	}

	response.Data(w, http.StatusOK, id)
}
