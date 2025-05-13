-- name: InsertAsset :one
INSERT INTO assets 
    (owner, price, detail)
VALUES 
    ($1, $2, $3)
RETURNING *;

-- name: GetAssetById :one
SELECT 
  a.id,
  a.owner,
  a.price,
  a.detail,
  a.status,
  a.created_at,
  a.updated_at,
  ac.id AS contact_id,
  ac.contact_name,
  ac.contact_detail,
  ai.id AS image_id,
  ai.image_url
FROM assets a
LEFT JOIN asset_contacts ac ON ac.asset_id = a.id
LEFT JOIN asset_images ai ON ai.asset_id = a.id
WHERE a.id = $1;


-- name: GetAssetsByUsername :many
SELECT 
  a.id,
  a.owner,
  a.price,
  a.detail,
  a.status,
  a.created_at,
  a.updated_at,
  ac.id AS contact_id,
  ac.contact_name,
  ac.contact_detail,
  ai.id AS image_id,
  ai.image_url
FROM assets a
LEFT JOIN asset_contacts ac ON ac.asset_id = a.id
LEFT JOIN asset_images ai ON ai.asset_id = a.id
WHERE a.owner = $1;

-- name: UpdateAsset :exec
UPDATE assets
SET price = $1, detail = $2
WHERE id = $3;

-- name: DeleteAsset :exec
DELETE FROM assets
WHERE id = $1;
