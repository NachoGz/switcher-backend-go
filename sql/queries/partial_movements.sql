-- name: CreatePartialMovement :one
INSERT INTO
	partial_movements (id, pos_from_x, pos_from_y, pos_to_x, pos_to_y, game_id, player_id, movement_card_id)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UndoMovement :exec
WITH last_movement AS (
	SELECT id FROM partial_movements pm
	WHERE pm.game_id = $1 AND pm.player_id = $2
	ORDER BY id DESC
	LIMIT 1
)
DELETE FROM partial_movements
WHERE id in (SELECT id from last_movement);

-- name: GetPartialMovementsByPlayer :many
SELECT * 
FROM partial_movements
WHERE game_id = $1 AND player_id = $2;

-- name: UndoMovementByID :exec
DELETE FROM partial_movements
WHERE id = $1;

-- name: DeleteAllPartialMovementsByPlayer :exec
DELETE FROM partial_movements
WHERE player_id = $1;