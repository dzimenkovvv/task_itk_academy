package server

import (
	"fmt"
	"net/http"
	"test_task_itk_academy/handler"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
}

func New(port string, h *handler.WalletHandler) *Server {
	router := mux.NewRouter()

	router.Path("/api/v1/wallet/create").Methods(http.MethodPost).HandlerFunc(h.CreateWallet)
	router.Path("/api/v1/wallet").Methods(http.MethodPost).HandlerFunc(h.Operation)
	router.Path("/api/v1/wallets/{WALLET_UUID}").Methods(http.MethodGet).HandlerFunc(h.GetBalance)

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
