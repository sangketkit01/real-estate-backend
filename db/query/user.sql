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
