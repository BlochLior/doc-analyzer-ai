-- name: CreateDocument :one
INSERT INTO documents (
    title,
    content,
    summary,
    ai_model
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetDocument :one
SELECT id, title, content, summary, ai_model, created_at, updated_at
FROM documents
WHERE id = $1;

-- name: ListDocuments :many
SELECT id, title, content, summary, ai_model, created_at, updated_at
FROM documents
ORDER BY created_at DESC;

-- name: UpdateDocument :one
UPDATE documents
SET
    title = $2,
    content = $3,
    summary = $4,
    ai_model = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1;