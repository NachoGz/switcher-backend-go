-- name: CreateGame :one
INSERT INTO games (id, name, max_players, min_players, is_private, password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAvailableGames :many
SELECT games.*
FROM games
JOIN game_state ON games.id = game_state.game_id
WHERE game_state.state = 'waiting';

-- name: GetGameById :one
SELECT *
FROM games
WHERE id = $1;