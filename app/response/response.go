// Package response marshall response to JSON and send it
// to the client with proper HTTP status.
package response

import (
	"encoding/json"
	"net/http"
)

type message struct {
	Success bool    `json:"success"`
	ErrCode string  `json:"err_code,omitempty"`
	ID      string  `json:"id,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
}

// WalletError status to the client.
func WalletError(w http.ResponseWriter, status int, err, id string) {
	msg := message{
		ErrCode: err,
		ID:      id,
	}

	sendJSON(w, status, msg)

}

// WalletSuccess status to the client.
func WalletSuccess(w http.ResponseWriter, status int, id string) {
	msg := message{
		Success: true,
		ID:      id,
	}

	sendJSON(w, status, msg)
}

// OperationError status to the client.
func OperationError(w http.ResponseWriter, status int, err string, amount float64) {
	msg := message{
		ErrCode: err,
		Amount:  amount,
	}

	sendJSON(w, status, msg)

}

// OperationSuccess status to the client.
func OperationSuccess(w http.ResponseWriter, status int, amount float64) {
	msg := message{
		Success: true,
		Amount:  amount,
	}

	sendJSON(w, status, msg)
}

// Data returns marshalled data to the client.
func Data(w http.ResponseWriter, status int, res any) {
	sendJSON(w, status, res)
}

func sendJSON(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
