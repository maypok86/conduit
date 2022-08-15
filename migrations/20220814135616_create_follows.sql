-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS follows (
    followee_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    follower_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CHECK (follower_id != followee_id),
    PRIMARY KEY (followee_id, follower_id)
);

CREATE INDEX IF NOT EXISTS follows_follower_id_idx ON follows (follower_id);
CREATE INDEX IF NOT EXISTS follows_followee_id_idx ON follows (followee_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS follows;
-- +goose StatementEnd
