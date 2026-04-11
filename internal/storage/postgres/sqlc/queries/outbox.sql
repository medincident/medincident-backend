-- name: AppendOutboxEvent :exec
INSERT INTO outbox.events (
    event_id,
    occurred_at,
    aggregate_type,
    aggregate_id,
    correlation_id,
    subject,
    headers,
    payload,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);
