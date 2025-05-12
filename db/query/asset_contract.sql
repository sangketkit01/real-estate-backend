-- name: InsertAssetContact :one
INSERT INTO asset_contacts
    (asset_id, contact_name, contact_detail)
VALUES
    ($1, $2, $3)
RETURNING *;
