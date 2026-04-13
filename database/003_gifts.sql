-- +goose Up
CREATE TABLE IF NOT EXISTS gifts(
    id BIGSERIAL PRIMARY KEY,
    wishlist_id int8 NOT NULL REFERENCES wishlists(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    link TEXT,
    priority int NOT NULL CHECK (priority BETWEEN 1 AND 5),
    booked bool NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE IF EXISTS gifts;