package api

import "wallet-service/internal/service"

type WalletHandler struct {
	service service.WalletService
}

func NewHandler(service service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}
