-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE chat (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE chat;
-- +goose StatementEnd
