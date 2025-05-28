-- name: CreateTag :one
INSERT INTO tags (name) VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: GetTagByName :one 
SELECT id, name, created_at
FROM tags
WHERE name = $1;

-- name: ListTags :many
SELECT id, name, created_at
FROM tags
ORDER BY name ASC;

-- name: GetTagsForDocument :many
SELECT tags.id, tags.name, tags.created_at
FROM tags
JOIN document_tags ON tags.id = document_tags.tag_id
WHERE document_tags.document_id = $1;

-- name: AddDocumentTag :exec
INSERT INTO document_tags (document_id, tag_id)
VALUES ($1, $2)
ON CONFLICT (document_id, tag_id) DO NOTHING;

-- name: RemoveDocumentTag :exec
DELETE FROM document_tags
WHERE document_id = $1 AND tag_id = $2;