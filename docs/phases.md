# Project Phases

This document outlines the development phases for the Digital Wallet API.

## Phase 1: Database & Core Logic (Completed ✅)
- [x] Design PostgreSQL relational schema (`users`, `accounts`, `transfers`).
- [x] Setup Docker Compose for isolated DB environment.
- [x] Write schema migrations using `golang-migrate`.
- [x] Write raw SQL queries for CRUD operations.
- [x] Generate type-safe Go code using `sqlc`.
- [x] Implement robust database transaction manager (`Store`) for atomic P2P transfers.

## Phase 2: Security & Utilities (Completed ✅)
- [x] Implement Password Hashing utilities.
- [x] Implement JWT generation and validation utilities.
- [x] Create Authentication middleware to protect API routes.
- [x] Create Logging middleware for request tracking.

## Phase 3: API Implementation (Completed ✅)
- [x] Setup Go 1.22+ native `net/http` server with routing.
- [x] Implement `Users` handlers (Registration, Login).
- [x] Implement `Accounts` handlers (Create account, Get balance, List accounts).
- [x] Implement `Transfers` handlers (Execute atomic transfer).
- [x] Integrate `Auth` and `Logger` middlewares with router.
- [x] Implement graceful shutdown for the HTTP server.

## Phase 4: Testing & Quality Assurance (Completed ✅)
- [x] Write unit tests for utilities (`password_test.go`, `jwt_test.go`, `config_test.go`).
- [x] Write unit tests for middlewares (`auth_test.go`, `logger_test.go`).
- [x] Write modular table-driven unit tests for API handlers (`user_test.go`, `account_test.go`, `transfer_test.go`, `server_test.go`).
- [x] Implement database deadlock integration tests (`store_test.go`).
- [x] Setup CI pipeline (`.github/workflows/ci.yml`) to run tests automatically.

## Phase 5: Production Readiness (In Progress 🚧)
- [x] Setup configuration management (`util/config.go` using `joho/godotenv`).
- [ ] Write Dockerfile for the Go API application.
- [ ] Finalize production deployment strategy.
