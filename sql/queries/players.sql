-- name: CreatePlayer :one
INSERT INTO players (id, name, turn, game_id, game_state_id, host)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CountPlayers :one
SELECT COUNT(*) FROM players WHERE game_id = $1;