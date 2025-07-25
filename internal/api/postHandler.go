package api

import (
	"encoding/json"
	"net/http"

	"wallet-service/internal/models"

	"github.com/google/uuid"
)

func (h *WalletHandler) PostWallet(w http.ResponseWriter, r *http.Request) {
	var op models.WalletOperation
	if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if op.WalletID == uuid.Nil {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	if op.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	wallet, err := h.service.ProcessOperation(r.Context(), op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}
