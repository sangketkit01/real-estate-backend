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
SET name = coalesce(sqlc.narg(name), name) , email = coalesce(sqlc.narg(email), email) , phone = coalesce(sqlc.narg(phone), phone) , 
    profile_url = coalesce(sqlc.narg(profile_url), profile_url)
WHERE  username = sqlc.arg(username);

-- name: GetUserPassword :one
SELECT password FROM users
WHERE username = $1;

-- name: UpdateUserPassword :exec
UPDATE users 
SET password = $1
WHERE username = $2;