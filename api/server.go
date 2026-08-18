package api

import (
	"net/http"

	"wallet-api/db/sqlc"
	"wallet-api/middleware"
	"wallet-api/util"
)

type Server struct {
	config util.Config
	store  db.Store
	router *http.ServeMux
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
		router: http.NewServeMux(),
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	server.router.HandleFunc("POST /users", server.createUser)
	server.router.HandleFunc("POST /users/login", server.loginUser)

	authMiddleware := middleware.Auth(server.config.JWTSecret)

	server.router.Handle("POST /accounts", authMiddleware(http.HandlerFunc(server.createAccount)))
	server.router.Handle("GET /accounts/{id}", authMiddleware(http.HandlerFunc(server.getAccount)))
	server.router.Handle("GET /accounts", authMiddleware(http.HandlerFunc(server.listAccounts)))
	server.router.Handle("POST /transfers", authMiddleware(http.HandlerFunc(server.createTransfer)))
}

func (server *Server) Router() http.Handler {
	return middleware.Logger(server.router)
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.Router())
}
