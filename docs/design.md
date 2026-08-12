# System & Database Design

## 1. Database Schema Design
The relational database is the core of the financial ledger. We use PostgreSQL.

### `users` Table
Stores authentication and profile data.
- `id`: BIGSERIAL (Primary Key)
- `email`: VARCHAR (Unique, Not Null)
- `password_hash`: VARCHAR (Not Null)
- `created_at`: TIMESTAMPTZ (Default NOW)

### `accounts` Table
A user can have multiple accounts (wallets), potentially in different currencies.
- `id`: BIGSERIAL (Primary Key)
- `user_id`: BIGINT (Foreign Key to `users.id`, ON DELETE CASCADE)
- `balance`: BIGINT (Not Null) - Stored in cents/smallest denomination to avoid floating-point errors.
- `currency`: VARCHAR(3) (Not Null)
- `created_at`: TIMESTAMPTZ (Default NOW)
- **Constraint**: `CHECK (balance >= 0)` - A balance can never drop below zero.

### `transfers` Table
Records all movements of money between accounts.
- `id`: BIGSERIAL (Primary Key)
- `from_account_id`: BIGINT (Foreign Key to `accounts.id`)
- `to_account_id`: BIGINT (Foreign Key to `accounts.id`)
- `amount`: BIGINT (Not Null, `CHECK (amount > 0)`)
- `created_at`: TIMESTAMPTZ (Default NOW)

## 2. API Design (RESTful)

### Users
- `POST /users`: Register a new user. Returns user info (excluding password).
- `POST /users/login`: Authenticate and receive a JWT token.

### Accounts
- `POST /accounts`: Create a new account for the authenticated user.
- `GET /accounts/{id}`: Retrieve account balance and details (Must belong to authenticated user).
- `GET /accounts`: List all accounts belonging to the authenticated user.

### Transfers
- `POST /transfers`: Execute a money transfer.
  - Payload: `{ "from_account_id": 1, "to_account_id": 2, "amount": 1000, "currency": "USD" }`
  - Validates that `from_account_id` belongs to the authenticated user.
  - Validates sufficient balance.
  - Executes as an atomic DB transaction.

## 3. Concurrency & Deadlock Prevention
To prevent database deadlocks during concurrent transfers (e.g., Account A sends to Account B, while Account B sends to Account A simultaneously):
- Queries use `SELECT ... FOR UPDATE` to lock rows during a transaction.
- When updating multiple accounts in a transfer transaction, we must lock accounts in a consistent order (e.g., always lock the account with the smaller ID first). This deterministic locking sequence prevents circular waits (deadlocks).
