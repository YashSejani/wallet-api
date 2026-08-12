# Architecture Overview

## 1. High-Level Architecture
This project strictly follows **Clean Architecture** principles to separate concerns, making the codebase testable, scalable, and easy to maintain. The application is divided into layers where dependencies point inwards.

### Layers:
1. **Routing & Delivery (Handlers)**: Uses Go 1.22+ native `net/http` multiplexer. Handles HTTP requests, parses JSON payloads, and returns HTTP responses.
2. **Business Logic (Services)**: Contains the core logic (e.g., validating transfers, orchestrating operations). Currently, some logic resides in the `Store` struct to manage DB transactions.
3. **Data Access (Repositories)**: Uses `sqlc` to auto-generate type-safe Go code from raw SQL queries. Interfaces with PostgreSQL.

## 2. Directory Structure
- `/db/migration/`: Contains `.up.sql` and `.down.sql` files managed by `golang-migrate`.
- `/db/query/`: Raw SQL files defining the database operations.
- `/db/sqlc/`: Auto-generated Go code by `sqlc` for database access, plus `store.go` for managing complex ACID transactions.
- `/middleware/`: HTTP middlewares such as `auth.go` (JWT verification) and `logger.go` (request logging).
- `/util/`: Helper functions like password hashing (`password.go`) and token generation (`jwt.go`).
- `/docs/`: Project documentation.

## 3. Database Architecture
- **PostgreSQL**: Chosen for its robust ACID compliance, which is critical for financial applications.
- **Transactions**: Complex operations like transferring money involve multiple steps (deducting from A, adding to B, recording the transfer). We use `pgxpool` and custom transaction blocks (`execTx`) to guarantee atomicity.

## 4. Middleware Pipeline
Requests pass through a standardized middleware chain:
`Request -> Logger Middleware -> Auth Middleware (if protected) -> Handler -> Response`
