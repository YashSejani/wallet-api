package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"wallet-api/db/sqlc"
	"wallet-api/middleware"
)

type createAccountRequest struct {
	Currency string `json:"currency"`
}

func (server *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("unauthorized")))
		return
	}

	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	if !IsSupportedCurrency(req.Currency) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("unsupported currency")))
		return
	}

	arg := db.CreateAccountParams{
		UserID:   userID,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(r.Context(), arg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(account)
}

func (server *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("unauthorized")))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("invalid account id")))
		return
	}

	account, err := server.store.GetAccount(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("account not found")))
		return
	}

	if account.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("account does not belong to the authenticated user")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(account)
}

type listAccountsRequest struct {
	PageID   int32 `json:"page_id"`
	PageSize int32 `json:"page_size"`
}

func (server *Server) listAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("unauthorized")))
		return
	}

	pageIDStr := r.URL.Query().Get("page_id")
	pageSizeStr := r.URL.Query().Get("page_size")

	pageID := int32(1)
	pageSize := int32(10)

	if val, err := strconv.Atoi(pageIDStr); err == nil && val > 0 {
		pageID = int32(val)
	}
	if val, err := strconv.Atoi(pageSizeStr); err == nil && val > 0 && val <= 100 {
		pageSize = int32(val)
	}

	arg := db.ListAccountsParams{
		UserID: userID,
		Limit:  pageSize,
		Offset: (pageID - 1) * pageSize,
	}

	accounts, err := server.store.ListAccounts(r.Context(), arg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(accounts)
}
