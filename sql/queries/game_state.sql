-- name: CreateGameState :one
INSERT INTO game_state (id, state, game_id, current_player_id, forbidden_color)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;