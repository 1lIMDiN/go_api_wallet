package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wallet-service/internal/models"

	"github.com/google/uuid"
)

// The GetWallet function outputs the wallet data
func (r *WalletRepo) GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.GetContext(ctx, &wallet, "SELECT id, balance FROM wallets WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	return &wallet, nil
}
