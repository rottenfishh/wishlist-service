-- +goose Up
ALTER TABLE wishlists
    ADD CONSTRAINT wishlist_token_unique UNIQUE(token);
-- +goose Down
ALTER TABLE wishlists
    DROP CONSTRAINT wishlist_token_unique;