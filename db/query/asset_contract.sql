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
SET contact_name = $1, contact_detail = $2
WHERE id = $3;


-- name: GetContact :one
SELECT * FROM asset_contacts 
WHERE id = $1;