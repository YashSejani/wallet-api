# Developer Rules & Guidelines

Welcome to the Digital Wallet API project! To ensure code quality, maintainability, and security, all developers must adhere to the following rules:

## 1. Database & SQL Rules
- **No Raw SQL in Go Code**: Do not write `db.Exec("SELECT...")` manually in Go files. All SQL queries must be written in `/db/query/*.sql`.
- **Use `sqlc`**: After modifying or adding queries in `/db/query/`, run `sqlc generate` to update the type-safe Go bindings in `/db/sqlc/`.
- **Migrations only via `golang-migrate`**: Never modify the database schema manually via a GUI or CLI client. Always create a new migration file (`migrate create -ext sql -dir db/migration -seq <name>`) and run it.

## 2. Security Best Practices
- **Never Commit Secrets**: Do not hardcode database URLs, JWT secrets, or API keys in the source code. Always use environment variables (`.env` file). Use `.env.example` to show what variables are required.
- **Password Handling**: Never log, print, or return user passwords in API responses. Always hash passwords before storing them.

## 3. Transaction Management
- **Use `Store.execTx` for Multiple DB Ops**: Any business logic that requires multiple database inserts/updates (like a money transfer) MUST be wrapped in a database transaction to ensure atomicity. Do not execute them sequentially without a transaction.

## 4. Code Structure & Style
- **Clean Architecture**: Handlers should only handle HTTP, parsing, and validating input. They should not contain complex business logic or raw database calls.
- **Error Handling**: Always check for errors. Do not ignore them (`_`). Return meaningful HTTP status codes (e.g., 400 Bad Request for validation errors, 500 Internal Server Error for DB failures).
- **Naming Conventions**: Follow standard Go naming conventions (e.g., camelCase for variables, PascalCase for exported types/functions).

## 5. Middleware
- Use the `Auth` middleware to protect any route that requires a logged-in user. Retrieve the `userID` from the context set by the middleware.
