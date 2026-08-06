-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateAccount :one
INSERT INTO accounts (user_id, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 FOR UPDATE;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + $2
WHERE id = $1
RETURNING *;