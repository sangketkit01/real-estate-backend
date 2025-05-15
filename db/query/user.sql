-- name: CreateUser :one
INSERT INTO users 
    (username, name, email, phone, password) 
VALUES(
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE username = $1;

-- name: LoginUser :one
SELECT username, password FROM users
WHERE username = $1;

-- name: UpdateUser :exec
UPDATE users
SET name = $1 , email = $2 , phone = $3 , profile_url = coalesce($4, profile_url)
WHERE  username = $5;

-- name: GetUserPassword :one
SELECT password FROM users
WHERE username = $1;

-- name: UpdateUserPassword :exec
UPDATE users 
SET password = $1
WHERE username = $2;