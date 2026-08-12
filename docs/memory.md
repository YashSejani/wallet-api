# Project Memory & Current Status

## Current Status: In Progress 🚧
The project is currently transitioning from a concept to a fully functional backend API. The data layer, utility layer, and middleware layers have been successfully established. The next major step is to tie these components together using an HTTP router and handlers.

## Implemented Functionalities
The following components have been verified as implemented in the codebase:

1. **Database & Migrations**: 
   - `000001_init_schema.up.sql` correctly defines the schemas for `users`, `accounts`, and `transfers` with appropriate constraints (e.g., positive balance checks).
   - Indexes are set up for performance.

2. **Data Access Layer (`sqlc`)**:
   - `ledger.sql` contains queries for creating users, fetching users by email, creating accounts, getting accounts `FOR UPDATE`, and updating account balances.
   - `sqlc` has generated the Go interfaces (`models.go`, `ledger.sql.go`) using `pgx/v5`.
   - `store.go` provides a sophisticated transaction wrapper (`execTx`) to handle operations that span multiple queries, guaranteeing ACID properties and supporting rollbacks on failure.

3. **Security & Utilities**:
   - `password.go`: Functions for hashing passwords and verifying passwords using `golang.org/x/crypto/bcrypt`.
   - `jwt.go`: Implementation of JSON Web Tokens using `github.com/golang-jwt/jwt/v5`, including specific `Claims` for `user_id`.

4. **Middleware**:
   - `auth.go`: HTTP middleware that intercepts requests, extracts the Bearer token, validates the JWT, and injects the `UserID` into the request Context.
   - `logger.go`: HTTP middleware that logs the incoming request method, path, response status code, and latency.

## Pending Tasks / Next Steps
- Write the HTTP Handler functions to process incoming JSON requests and map them to the `Store` methods.
- Initialize the routing mechanism (`main.go`) to define paths (e.g., `POST /users`, `POST /transfers`) and apply the `Logger` and `Auth` middlewares appropriately.
- Connect the application to read environment variables (using a package like `viper` or standard `os.Getenv`) for the database connection string and JWT secret.
