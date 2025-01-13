-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE Email = $1;

-- name: UpdateUser :one
UPDATE users
SET updated_at = now(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: UpgradeUserToChirpyRed :one
UPDATE users
SET updated_at = now(), is_chirpy_red = true
WHERE id = $1
RETURNING *;