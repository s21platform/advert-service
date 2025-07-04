-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS advert_text
(
    id           SERIAL PRIMARY KEY,
    owner_uuid   UUID NOT NULL,
    title        TEXT NOT NULL,
    text_content TEXT,
    filter       JSONB,
    expired_at   TIMESTAMP,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_canceled BOOLEAN DEFAULT FALSE,
    canceled_at TIMESTAMP,
    is_banned BOOLEAN DEFAULT FALSE,
    banned_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS advert_text;
-- +goose StatementEnd
