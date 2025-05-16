-- name: InsertAssetContact :one
INSERT INTO asset_contacts
    (asset_id, contact_name, contact_detail)
VALUES
    ($1, $2, $3)
RETURNING *;


-- name: RemoveContact :exec
DELETE FROM asset_contacts 
WHERE id = $1;


-- name: UpdateContact :exec
UPDATE asset_contacts
SET contact_name = coalesce($1, contact_name), contact_detail = coalesce($2, contact_detail)
WHERE id = $3;


-- name: GetContact :one
SELECT * FROM asset_contacts 
WHERE id = $1;

-- name: GetAssetContacts :many
SELECT * FROM asset_contacts
WHERE asset_id = $1;