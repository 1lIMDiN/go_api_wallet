package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// The ProcessOperation function performs an operation to change the balance
func (r *WalletRepo) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentBalance int64
	err = tx.GetContext(ctx, &currentBalance, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", id)
	if err != nil {
		return fmt.Errorf("failed to get wallet balance: %w", err)
	}

	newBalance := currentBalance + amount
	if newBalance < 0 {
		return fmt.Errorf("insufficient funds")
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, id)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	return tx.Commit()
}
