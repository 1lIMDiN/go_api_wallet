package repository

import (
	"context"
	"fmt"

	"wallet-service/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) error
	CreateWallet(ctx context.Context, id uuid.UUID) error
}

type WalletRepo struct {
	db *sqlx.DB
}

func NewWallet(db *sqlx.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) CreateWallet(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO wallets (id, balance) VALUES ($1, 0)", id)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}
