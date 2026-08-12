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

## Phase 3: API Implementation (In Progress 🚧)
- [ ] Setup Go 1.22+ native `net/http` server.
- [ ] Implement `Users` handlers (Registration, Login).
- [ ] Implement `Accounts` handlers (Create account, Get balance).
- [ ] Implement `Transfers` handlers (Execute transfer).
- [ ] Integrate middlewares with routers.
- [ ] Implement graceful shutdown for the HTTP server.

## Phase 4: Testing & Quality Assurance (Upcoming 📅)
- [ ] Write unit tests for database CRUD operations (mocking/test DB).
- [ ] Write unit tests for API handlers using `httptest`.
- [ ] Implement database deadlock testing for concurrent transactions.
- [ ] Setup CI pipeline (e.g., GitHub Actions) to run tests automatically.

## Phase 5: Production Readiness (Upcoming 📅)
- [ ] Write Dockerfile for the Go API application.
- [ ] Setup configuration management (loading from env vars/files).
- [ ] Finalize production deployment strategy.
