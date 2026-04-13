-- +goose Up
CREATE TABLE IF NOT EXISTS wishlists(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ,
    token UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS wishlists;