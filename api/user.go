package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"wallet-api/db/sqlc"
	"wallet-api/util"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("email is required")))
		return
	}
	if len(req.Password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("password must be at least 6 characters")))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	user, err := server.store.CreateUser(r.Context(), arg)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(errorResponse(errors.New("email already registered")))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(rsp)
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("invalid email or password")))
		return
	}

	err = util.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse(errors.New("invalid email or password")))
		return
	}

	accessToken, err := util.GenerateToken(user.ID, server.config.JWTSecret, server.config.TokenDuration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rsp)
}
