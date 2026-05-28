package repositories

import (
	"context"
	"database/sql"
	"fmt"
)

// BidanRepository handles looking up bidan IDs from user IDs.
type BidanRepository struct {
	db *sql.DB
}

// NewBidanRepository creates a new BidanRepository instance.
func NewBidanRepository(db *sql.DB) *BidanRepository {
	return &BidanRepository{db: db}
}

// FindBidanIDByUserID returns the bidan.id for a given user_id.
func (r *BidanRepository) FindBidanIDByUserID(ctx context.Context, userID int) (int, error) {
	var bidanID int
	err := r.db.QueryRowContext(ctx, `SELECT id FROM bidan WHERE user_id = ?`, userID).Scan(&bidanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("bidan tidak ditemukan untuk user_id %d", userID)
		}
		return 0, fmt.Errorf("query bidan: %w", err)
	}
	return bidanID, nil
}
