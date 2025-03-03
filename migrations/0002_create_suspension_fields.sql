-- +goose Up
-- +goose StatementBegin
ALTER TABLE advert_text
    ADD COLUMN IF NOT EXISTS is_canceled BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS canceled_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS is_banned BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS banned_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE advert_text
    DROP COLUMN IF EXISTS is_canceled,
    DROP COLUMN IF EXISTS canceled_at,
    DROP COLUMN IF EXISTS is_banned,
    DROP COLUMN IF EXISTS banned_at;
-- +goose StatementEnd
