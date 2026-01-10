-- name: AddUser :exec
INSERT INTO users (
    access_token,
    showed
) VALUES (?, ?);

-- name: GetAccessToken :one
SELECT access_token FROM users
WHERE id = 1;

-- name: GetAPIKEY :one
SELECT api_key FROM users
WHERE
    id = 1
    AND showed = 0;

-- name: UpdateToken :exec
UPDATE users
SET access_token = ?
WHERE id = 1 ;

-- name: UpdateAPIKEY :exec
UPDATE users
SET api_key = ?
WHERE id = 1 ;
