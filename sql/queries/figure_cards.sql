-- name: CreateFigureCard :one
INSERT INTO
	figure_cards (id, show, player_id, game_id, type, blocked, soft_blocked)
VALUES
	($1, $2, $3, $4, $5, $6 ,$7)
RETURNING *;

