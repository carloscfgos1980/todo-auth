-- name: CreateTodo :one
INSERT INTO todos (created_at, updated_at, title, completed, user_id)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;