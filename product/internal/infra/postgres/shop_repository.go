package postgres

import (
	"database/sql"
	"github.com/google/uuid"
)

type shopRepository struct {
	db *sql.DB
}

func NewShopRepository(db *sql.DB) *shopRepository {
	return &shopRepository{db: db}
}

// IsShopOwner checks if the given userID is the owner of the shopID.
// This directly queries the 'shops' table.
func (r *shopRepository) IsShopOwner(userIDStr, shopIDStr string) (bool, error) {
	var ownerID uuid.UUID

	// Ensure string IDs are valid UUIDs before querying if necessary, or let DB handle type error
	// For this implementation, we assume shopIDStr is a valid UUID string.
	// userIDStr should also be a valid UUID string representing the customer's ID.

	query := `SELECT owner_id FROM shops WHERE id = $1`
	err := r.db.QueryRow(query, shopIDStr).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Shop not found, so user cannot be the owner
			return false, nil // Or return an error like "shop not found"
		}
		return false, err
	}

	return ownerID.String() == userIDStr, nil
}
