package domain

import (
	"time"

	"github.com/google/uuid"
)

// Shop represents a shop entity.
type Shop struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
