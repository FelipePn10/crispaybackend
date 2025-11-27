-- name: CreateVerificationSession :one
INSERT INTO verification_sessions (
    user_id,
    session_id,
    didit_session_id,
    verification_url,
    user_email,
    user_first_name,
    user_last_name,
    status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetVerificationSessionByID :one
SELECT * FROM verification_sessions 
WHERE id = $1 LIMIT 1;

-- name: GetVerificationSessionBySessionID :one
SELECT * FROM verification_sessions 
WHERE session_id = $1 LIMIT 1;

-- name: GetVerificationSessionByDiditSessionID :one
SELECT * FROM verification_sessions 
WHERE didit_session_id = $1 LIMIT 1;

-- name: UpdateVerificationSessionStatus :one
UPDATE verification_sessions 
SET 
    status = $2,
    updated_at = NOW(),
    completed_at = CASE 
        WHEN $2 IN ('approved', 'rejected', 'failed') THEN NOW() 
        ELSE completed_at 
    END
WHERE session_id = $1
RETURNING *;

-- name: UpdateDiditSessionData :one
UPDATE verification_sessions 
SET 
    didit_session_id = $2,
    verification_url = $3,
    updated_at = NOW()
WHERE session_id = $1
RETURNING *;

-- name: ListVerificationSessionsByUserID :many
SELECT * FROM verification_sessions 
WHERE user_id = $1 
ORDER BY created_at DESC;

-- name: ListVerificationSessionsByStatus :many
SELECT * FROM verification_sessions 
WHERE status = $1 
ORDER BY created_at DESC;

-- name: CreateWebhookEvent :one
INSERT INTO webhook_events (
    event_type,
    session_id,
    payload
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetWebhookEventsBySessionID :many
SELECT * FROM webhook_events 
WHERE session_id = $1 
ORDER BY created_at DESC;