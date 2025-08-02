-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    unique_url text unique not null
);
CREATE TABLE user_links(
     id SERIAL PRIMARY KEY,
    user_id integer not null REFERENCES users(id),
    urls_id integer not null REFERENCES urls(id),
     create_at timestamp(0) DEFAULT now()
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
DROP TABLE user_links;
-- +goose StatementEnd
