package storage

import (
	"fmt"
	"time"
)

type ID int64

var BadID = ID(-1)

type Element struct {
	ID      ID        `json:"id,omitempty"`
	Payload string    `json:"payload,omitempty"`
	Added   time.Time `json:"added,omitempty"`
}

func (e Element) String() string {
	return fmt.Sprintf("ID: %v\nPayload: %v\nTimestamp: %v\n", e.ID, e.Payload, e.Added)
}
