-- +goose Up
CREATE TABLE document_tags (
    document_id UUID NOT NULL references documents(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL references tags(id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, tag_id)
);

-- +goose Down
DROP TABLE document_tags;