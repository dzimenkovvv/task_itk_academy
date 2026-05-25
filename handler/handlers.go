package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"test_task_itk_academy/database"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type WalletHandler struct {
	repo *database.WalletRepository
}

func NewWalletHandler(repo *database.WalletRepository) *WalletHandler {
	return &WalletHandler{
		repo: repo,
	}
}

type operationRequest struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        float64   `json:"amount"`
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallet, err := h.repo.Create(r.Context())
	if err != nil {
		writeError(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, wallet)
}

func (h *WalletHandler) Operation(w http.ResponseWriter, r *http.Request) {
	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		writeError(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	var (
		wallet *database.Wallet
		err    error
	)

	switch req.OperationType {
	case "DEPOSIT":
		wallet, err = h.repo.Deposit(r.Context(), req.WalletID, req.Amount)
	case "WITHDRAW":
		wallet, err = h.repo.Withdraw(r.Context(), req.WalletID, req.Amount)
	default:
		writeError(w, "OperationType must be DEPOSIT or WITHDRAW", http.StatusBadRequest)
	}

	if err != nil {
		switch {
		case errors.Is(err, database.ErrWalletNotFound):
			writeError(w, "Wallet not found", http.StatusNotFound)
		case errors.Is(err, database.ErrInsufficientFunds):
			writeError(w, "Insufficient funds", http.StatusUnprocessableEntity)
		default:
			writeError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	rawID := mux.Vars(r)["WALLET_UUID"]

	walletID, err := uuid.Parse(rawID)
	if err != nil {
		writeError(w, "Invalid wallet id", http.StatusBadRequest)
	}

	wallet, err := h.repo.GetBalance(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, database.ErrWalletNotFound) {
			writeError(w, "Wallet not found", http.StatusNotFound)
			return
		}
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, wallet)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, status, errorResponse{Error: msg})
}
