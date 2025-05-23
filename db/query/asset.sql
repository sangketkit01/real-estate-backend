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
  MIN(ac.id) AS contact_id,
  MIN(ac.contact_name) AS contact_name,
  MIN(ac.contact_detail) AS contact_detail,
  MIN(ai.id) AS image_id,
  MIN(ai.image_url) AS image_url
FROM assets a
LEFT JOIN asset_contacts ac ON ac.asset_id = a.id
LEFT JOIN asset_images ai ON ai.asset_id = a.id
WHERE a.owner = $1
GROUP BY a.id
ORDER BY a.id DESC
LIMIT $2 OFFSET $3;



-- name: GetAllAssets :many
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
ORDER BY a.id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAsset :exec
UPDATE assets
SET price = coalesce($1, price), detail = coalesce($2, detail)
WHERE id = $3;

-- name: DeleteAsset :exec
DELETE FROM assets
WHERE id = $1;

-- name: GetAssetCount :one
SELECT count(id) FROM assets;

-- name: GetAssetCountByUsername :one
SELECT count(id) FROM assets
WHERE owner = $1;