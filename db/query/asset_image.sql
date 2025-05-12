-- name: InsertAssetImage :one
INSERT INTO asset_images 
    (asset_id, image_url)
VALUES 
    ($1, $2)
RETURNING *;