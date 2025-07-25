package service

import (
	"context"
	"fmt"

	"wallet-service/internal/models"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
)

type WalletService interface {
	GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	ProcessOperation(ctx context.Context, op models.WalletOperation) (*models.Wallet, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	return s.repo.GetWallet(ctx, id)
}

func (s *walletService) ProcessOperation(ctx context.Context, op models.WalletOperation) (*models.Wallet, error) {
	var amount int64
	switch op.OperationType {
	case models.Deposit:
		if op.Amount <= 0 {
			return nil, fmt.Errorf("deposit amount must be positive")
		}
		amount = op.Amount
	case models.Withdraw:
		if op.Amount <= 0 {
			return nil, fmt.Errorf("withdraw amount must be positive")
		}
		amount = -op.Amount
	default:
		return nil, fmt.Errorf("invalid operation type")
	}

	// Проверяем существует ли кошелек, если нет - создаем
	_, err := s.repo.GetWallet(ctx, op.WalletID)
	if err != nil {
		if err := s.repo.CreateWallet(ctx, op.WalletID); err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}
	}

	if err := s.repo.UpdateBalance(ctx, op.WalletID, amount); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	return s.repo.GetWallet(ctx, op.WalletID)
}
