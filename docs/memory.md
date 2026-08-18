# Project Memory & Current Status

## Current Status: Phase 4 Completed ✅ | Phase 5 In Progress 🚧
The core API server, database layer, transaction engine, security utilities, HTTP handlers, modular unit test suite, and GitHub Actions CI pipeline have been fully implemented and verified.

## Implemented Functionalities

1. **Database & Migrations**: 
   - `000001_init_schema.up.sql` defines schemas for `users`, `accounts`, and `transfers` with constraints (e.g., balance >= 0, amount > 0).
   - Foreign keys and indexes configured for query performance.

2. **Data Access Layer (`sqlc`) & Transaction Engine**:
   - `ledger.sql` defines queries for creating users, fetching users by email, creating accounts, fetching accounts (`GetAccount`, `GetAccountForUpdate`, `ListAccounts`), updating balances, and recording transfers.
   - `sqlc` generates type-safe Go bindings (`models.go`, `ledger.sql.go`) using `pgx/v5`.
   - `store.go` provides a `Store` interface and `SQLStore` implementation supporting `TransferTx` for atomic P2P transfers with deterministic row locking (smaller account ID updated first) to prevent database deadlocks.

3. **Security & Utilities**:
   - `password.go`: Password hashing and verification using `bcrypt`.
   - `jwt.go`: Signed JWT token creation and validation with standard `Claims`.
   - `config.go`: Environment configuration loader using `joho/godotenv` loading `.env` variables (`DB_SOURCE`, `SERVER_ADDRESS`, `JWT_SECRET`, `JWT_DURATION`).

4. **Middleware**:
   - `auth.go`: Intercepts protected requests, verifies Bearer tokens, and injects `UserID` into request Context.
   - `logger.go`: Logs HTTP method, URL path, response status code, and latency.

5. **API Layer (`/api`)**:
   - `user.go`: `POST /users` (registration) & `POST /users/login` (authentication, returns JWT).
   - `account.go`: `POST /accounts` (wallet creation), `GET /accounts/{id}` (authorization check), `GET /accounts` (paginated list).
   - `transfer.go`: `POST /transfers` (validates account ownership, currency match, sufficient balance, and executes `TransferTx`).
   - `server.go`: Configures Go 1.22+ `net/http` router with middleware chains.
   - `validator.go`: Input validation & standardized error response formatting.

6. **Application Entrypoint (`main.go`)**:
   - Connects to PostgreSQL via `pgxpool`, runs HTTP server, and implements graceful shutdown on `SIGINT` / `SIGTERM`.

7. **Testing & Quality Assurance (`Phase 4`)**:
   - `util/password_test.go`: Tests bcrypt hashing, verification, mismatch, and salt uniqueness.
   - `util/jwt_test.go`: Tests token generation, validation, expiration, and invalid secret signature handling.
   - `middleware/auth_test.go`: Tests Auth middleware with valid tokens, missing headers, malformed headers, and expired tokens.
   - `middleware/logger_test.go`: Tests status code capturing and request execution logging.
   - `api/mock_store_test.go`: Shared MockStore implementing `db.Store` interface for handler unit testing.
   - `api/user_test.go`: Table-driven tests for user registration and authentication endpoints.
   - `api/account_test.go`: Table-driven tests for wallet account creation, retrieval, and listing.
   - `api/transfer_test.go`: Table-driven tests for P2P transfers covering ownership, currency validation, balance checks, and destination validation.
   - `api/server_test.go`: Server router setup test.
   - `db/sqlc/store_test.go`: Concurrency and deadlock testing for `TransferTx`.
   - `.github/workflows/ci.yml`: GitHub Actions workflow automatically testing and building code on git push/pull requests.

## Pending Tasks / Next Steps (Phase 5)
- Write a multi-stage `Dockerfile` for the Go application.
- Setup production deployment configuration and instructions.
