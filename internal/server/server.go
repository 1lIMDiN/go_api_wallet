package server

import (
	"fmt"
	"net/http"

	"wallet-service/internal/api"
)

func Run(port string, handler *api.WalletHandler) error {
	api.Init(handler)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), api.Rout)
}
