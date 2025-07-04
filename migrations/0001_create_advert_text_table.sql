-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS adverts
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_uuid   UUID NOT NULL,
    title        TEXT NOT NULL,
    text_content TEXT,
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
DROP TABLE IF EXISTS adverts;
DROP EXTENSION IF EXISTS pgcrypto;
-- +goose StatementEnd
