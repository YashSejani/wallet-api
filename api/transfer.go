package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"wallet-api/db/sqlc"
	"wallet-api/middleware"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

func (server *Server) createTransfer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("unauthorized")))
		return
	}

	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	if req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("amount must be greater than zero")))
		return
	}

	if req.FromAccountID == req.ToAccountID {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("from_account_id and to_account_id cannot be the same")))
		return
	}

	if !IsSupportedCurrency(req.Currency) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("unsupported currency")))
		return
	}

	// Validate source account ownership, currency match, and sufficient balance
	fromAccount, err := server.store.GetAccount(r.Context(), req.FromAccountID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("source account not found")))
		return
	}

	if fromAccount.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("source account does not belong to authenticated user")))
		return
	}

	if fromAccount.Currency != req.Currency {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(fmt.Errorf("source account currency mismatch: %s vs %s", fromAccount.Currency, req.Currency)))
		return
	}

	if fromAccount.Balance < req.Amount {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("insufficient balance")))
		return
	}

	// Validate destination account existence and currency match
	toAccount, err := server.store.GetAccount(r.Context(), req.ToAccountID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("destination account not found")))
		return
	}

	if toAccount.Currency != req.Currency {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(fmt.Errorf("destination account currency mismatch: %s vs %s", toAccount.Currency, req.Currency)))
		return
	}

	// Execute atomic P2P transfer transaction
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(r.Context(), arg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
