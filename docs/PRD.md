# Product Requirements Document (PRD)
## Financial Ledger / Digital Wallet API

### 1. Project Goal
The primary objective of this project is to build a robust, secure, and persistent HTTP server for a digital wallet application. This project marks a transition from consuming external APIs to constructing a production-grade backend service capable of managing financial ledgers securely.

### 2. Target Audience
- Frontend developers building the user interface.
- Mobile developers integrating the wallet features.
- System administrators maintaining the deployment.

### 3. Core Features
- **User Management**: Secure user registration and authentication (login).
- **Authentication**: JWT-based session management for protecting API endpoints.
- **Account Management**: Users can create wallets/accounts in specific currencies, and check their balances.
- **P2P Money Transfers**: Users can transfer funds between accounts atomically.

### 4. Non-Functional Requirements
- **Security**: 
  - Passwords must be securely hashed (e.g., using bcrypt).
  - All sensitive endpoints must be protected by JWT middleware.
- **Data Integrity & Consistency**: 
  - All financial transactions (transfers) must use ACID-compliant database transactions to ensure balances are never in an invalid state.
  - Accounts cannot have negative balances.
- **Reliability & Scalability**: 
  - Use efficient connection pooling for database operations.
  - Adhere to Clean Architecture principles for maintainability.

### 5. Technology Stack
- **Language**: Go (Golang) 1.22+
- **Database**: PostgreSQL
- **Database Tools**: 
  - `sqlc` for type-safe database query generation.
  - `golang-migrate` for handling schema migrations.
- **Routing**: Go's native `net/http` (with 1.22 routing features) or `go-chi/chi`.
- **Infrastructure**: Docker & Docker Compose for local database spin-ups.
