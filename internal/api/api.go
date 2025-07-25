package api

import (
	"github.com/gorilla/mux"
)

var Rout = mux.NewRouter()

func Init(h *WalletHandler) {
	Rout.HandleFunc("/api/v1/wallet", h.PostWallet).Methods("POST")
	Rout.HandleFunc("/api/v1/wallets/{walletId}", h.GetWallet).Methods("GET")
}
