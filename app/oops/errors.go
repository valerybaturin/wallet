// Package oops is an error package for our app.
package oops

import "errors"

const (
	// ErrNoDataMessage - store has no data.
	ErrNoDataMessage = "no data"
	// ErrIntServMessage - internal server error
	ErrIntServMessage = "internal error"
	// ErrNotFoundMessage - requested wallet not found.
	ErrNotFoundMessage = "wallet not found"
	// ErrBadReqMessage - client error message.
	ErrBadReqMessage = "bad request"
	// ErrNotEnoMonMessage - not enough memory
	ErrNotEnoMonMessage = "not enough money"
)

// ErrNotFound — wallet not found.
var (
	ErrNotFound  = errors.New(ErrNotFoundMessage)
	ErrNotEnoMon = errors.New(ErrNotEnoMonMessage)
)
