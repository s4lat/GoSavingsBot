-- +goose Up
CREATE TABLE users (
   id BIGSERIAL PRIMARY KEY,
   username TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
