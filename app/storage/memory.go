// Package storage is a data storage.
package storage

// Wallet contains data fields for map.
type Wallet struct {
	Name, Status string
	Balance      float64
}

// Memory has map data to store wallets.
type Memory struct {
	Data map[string]Wallet
}

// NewMemory is a constructor which initiates a new map.
func NewMemory() *Memory {
	data := make(map[string]Wallet)
	return &Memory{
		Data: data,
	}
}
