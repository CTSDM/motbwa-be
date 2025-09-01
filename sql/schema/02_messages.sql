-- +goose Up
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    sent_at TIMESTAMP NOT NULL,
    content BYTEA,
    CONSTRAINT fk_sender_id FOREIGN KEY (sender_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    CONSTRAINT fk_receiver_id FOREIGN KEY (receiver_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE messages;
