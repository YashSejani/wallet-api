# Architecture Overview

## 1. High-Level Architecture
This project strictly follows **Clean Architecture** principles to separate concerns, making the codebase testable, scalable, and easy to maintain. The application is divided into layers where dependencies point inwards.

### Layers:
1. **Routing & Delivery (`/api`)**: Uses Go 1.22+ native `net/http` multiplexer. Handles HTTP requests, parses JSON payloads, validates input, and formats JSON HTTP responses.
2. **Business Logic & Transactions (`/db/sqlc/store.go`)**: Contains core financial logic (validating transfers, orchestrating ACID database transactions via `TransferTx`).
3. **Data Access (`/db/sqlc`)**: Uses `sqlc` to auto-generate type-safe Go code from raw SQL queries. Interfaces with PostgreSQL using `pgx/v5`.
4. **Configuration & Utilities (`/util`, `/middleware`)**: Configuration loader (`config.go`), password hashing (`password.go`), JWT tokens (`jwt.go`), and HTTP middlewares (`auth.go`, `logger.go`).

## 2. Directory Structure
- `/api/`: REST handlers (`user.go`, `account.go`, `transfer.go`), routing (`server.go`), input validation (`validator.go`), and handler unit tests (`main_test.go`).
- `/db/migration/`: `.up.sql` and `.down.sql` files managed by `golang-migrate`.
- `/db/query/`: Raw SQL files (`ledger.sql`) defining database queries.
- `/db/sqlc/`: Auto-generated Go code by `sqlc`, plus `store.go` implementing the `Store` interface and atomic `TransferTx` transactions.
- `/middleware/`: HTTP middlewares (`auth.go` for JWT verification, `logger.go` for request logging).
- `/util/`: Helper functions (`password.go`, `jwt.go`, `config.go`).
- `/docs/`: Project documentation and architecture specs.
- `/main.go`: Application entrypoint with DB pool connection and graceful server shutdown.

## 3. Database Architecture
- **PostgreSQL**: Chosen for its robust ACID compliance, critical for financial ledgers.
- **Transactions**: Atomic operations like P2P money transfers use `pgxpool` and custom transaction blocks (`execTx`).
- **Deadlock Prevention**: Deterministic row locking order (always updating the smaller account ID first) prevents circular locks during concurrent transfers.

## 4. Middleware Pipeline
Requests pass through a standardized middleware chain:
`Request -> Logger Middleware -> Auth Middleware (for protected routes) -> Handler -> Response`
