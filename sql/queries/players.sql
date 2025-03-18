-- name: CreatePlayer :one
INSERT INTO players (id, name, turn, game_id, game_state_id, host)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CountPlayers :one
SELECT COUNT(*) FROM players WHERE game_id = $1;

-- name: GetPlayersInGame :many
SELECT *
FROM players
WHERE game_id=$1;

-- name: AssignTurnPlayer :exec
UPDATE players
SET turn=$2, updated_at = NOW()
WHERE id=$1;

-- name: GetPlayerByID :one
SELECT *
FROM players
WHERE game_id=$1 AND id=$2;
