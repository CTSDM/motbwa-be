-- name: CreateMessage :one
INSERT INTO messages(id, sender_id, receiver_id, sent_at, content)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4)
RETURNING *;

-- name: GetMessagesBySender :many
SELECT * FROM messages
WHERE sender_id = $1;

-- name: GetMessagesByReceiver :many
SELECT * FROM messages
WHERE receiver_id = $1;

-- name: DeleteMessages :exec
DELETE FROM messages;
