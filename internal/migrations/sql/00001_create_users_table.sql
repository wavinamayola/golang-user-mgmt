-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT '', 
    age VARCHAR(255),
    created TIMESTAMP DEFAULT NOW(),
    updated TIMESTAMP DEFAULT NOW()
);

-- +goose StatementBegin
CREATE TRIGGER update_timestamp
BEFORE UPDATE
ON users
FOR EACH ROW
SET NEW.updated = NOW();
-- +goose StatementEnd

-- +goose Down
DROP TABLE users;
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_timestamp;
-- +goose StatementEnd
