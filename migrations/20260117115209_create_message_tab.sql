-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE message (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL REFERENCES chat(id) ON DELETE CASCADE,
    text VARCHAR(5000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE message;
-- +goose StatementEnd
