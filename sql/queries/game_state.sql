-- name: CreateGameState :one
INSERT INTO game_state (id, state, game_id, current_player_id, forbidden_color)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateGameState :exec
UPDATE game_state
SET state=$2
WHERE game_id=$1;

-- name: UpdateCurrentPlayer :exec
UPDATE game_state
SET current_player_id=$2
WHERE game_id=$1;

-- name: GetGameStateByGameID :one
SELECT *
FROM game_state
WHERE game_id=$1;