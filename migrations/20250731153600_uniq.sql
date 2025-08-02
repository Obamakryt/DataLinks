-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_links
ADD CONSTRAINT unique_user_urls UNIQUE (user_id, urls_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_links
DROP CONSTRAINT unique_user_urls;
-- +goose StatementEnd
