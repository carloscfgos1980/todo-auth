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

-- name: GetTodosByUserID :many
SELECT * FROM todos WHERE user_id = $1;

-- name: GetTodoByID :one
SELECT * FROM todos WHERE id = $1;